// Package archive implements the revision workflow: freeze the current version
// of a page under <name>.revisions/<label>/, publish it at a stable URL, and
// record the label -> URL list in front matter so the theme can render a
// version selector without guessing any URLs.
package archive

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ifeitao/hugo-theme-resound/tools/internal/fm"
)

// Schema is the front matter contract version. The theme refuses to render a
// list it does not understand, so a stale theme degrades instead of producing
// dead links.
const Schema = 1

// Options configures a revision.
type Options struct {
	// BaseURL is the page's current URL ("/posts/foo/"). When empty the tool
	// asks Hugo (`hugo list all`) and falls back to a path heuristic.
	BaseURL string
	// Label overrides the new version label. Empty means today's date.
	Label string
	// Stdout/Stderr receive progress and warnings. Defaults to os.Stdout/Stderr.
	Stdout io.Writer
	Stderr io.Writer
}

// Result describes what a run did.
type Result struct {
	File       string
	IsBundle   bool
	Archived   string // path of the archived copy (file or bundle dir)
	Label      string // label of the version that got archived
	NewLabel   string // label of the live version after this run
	BaseURL    string
	BaseSource string // "flag", "front matter", "hugo list all" or "path heuristic"
	Entries    []fm.Revision
}

type paths struct {
	root     string // Hugo project root
	file     string // the .md to edit
	dir      string // directory holding file
	name     string // page name (file base without .md, or bundle dir name)
	isBundle bool
	revDir   string // <dir>/<name>.revisions
	relPath  string // posix path relative to content/
}

// Revise archives the current version of target and promotes a new label.
func Revise(target string, o Options) (*Result, error) {
	stdout, stderr := writers(o)
	p, err := resolve(target)
	if err != nil {
		return nil, err
	}

	src, err := os.ReadFile(p.file)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", p.file, err)
	}
	doc, err := fm.Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p.file, err)
	}

	entries := doc.Revisions()
	curLabel := ""
	if n := len(entries); n > 0 {
		curLabel = entries[n-1].Label
	}
	if curLabel == "" {
		curLabel = datePart(doc.Date())
		if curLabel == "" {
			return nil, fmt.Errorf("%s: no date in front matter; add `date: YYYY-MM-DD` (it becomes the first version label)", p.file)
		}
	}

	newLabel := strings.TrimSpace(o.Label)
	if newLabel == "" {
		newLabel = time.Now().Format("2006-01-02")
	}
	if newLabel == curLabel {
		return nil, fmt.Errorf("version %q already exists; pass --label with a later date (e.g. --label %s-2) or wait a day", newLabel, newLabel)
	}
	if len(entries) > 0 && newLabel < curLabel {
		return nil, fmt.Errorf("label %q is older than the current version %q; labels must sort chronologically", newLabel, curLabel)
	}

	base, source, err := resolveBase(p, entries, o.BaseURL, stdout, stderr)
	if err != nil {
		return nil, err
	}
	archiveURL := base + "revisions/" + curLabel + "/"
	merged := mergeEntries(entries, curLabel, archiveURL, newLabel, base)

	if err := os.MkdirAll(p.revDir, 0o755); err != nil {
		return nil, fmt.Errorf("create %s: %w", p.revDir, err)
	}

	// 1. Freeze the current version.
	archived, err := writeArchive(p, src, curLabel, archiveURL, merged)
	if err != nil {
		return nil, err
	}

	// 2. Every archive carries the same, up-to-date list.
	if err := syncArchives(p.revDir, merged); err != nil {
		return nil, err
	}

	// 3. Promote the live page: new list + lastmod. `date` is deliberately left
	//    untouched — with date-based permalinks rewriting it moves the live URL
	//    and silently breaks every inbound link.
	live, err := fm.Parse(string(src))
	if err != nil {
		return nil, err
	}
	live.SetSchema(Schema)
	live.SetCurrentLabel(newLabel)
	live.SetRevisions(merged)
	live.SetLastmod(time.Now().Format("2006-01-02T15:04:05-07:00"))
	if err := os.WriteFile(p.file, live.Bytes(), 0o644); err != nil {
		return nil, fmt.Errorf("write %s: %w", p.file, err)
	}

	// 4. Remember how to undo it.
	if err := saveUndo(p.root, p.file, archived, string(src)); err != nil {
		fmt.Fprintf(stderr, "warning: could not save undo state: %v\n", err)
	}

	return &Result{
		File:       p.file,
		IsBundle:   p.isBundle,
		Archived:   archived,
		Label:      curLabel,
		NewLabel:   newLabel,
		BaseURL:    base,
		BaseSource: source,
		Entries:    merged,
	}, nil
}

// Undo reverts the last revision of target: restores the live page byte for
// byte and deletes the archived copy.
func Undo(target string) (*Result, error) {
	p, err := resolve(target)
	if err != nil {
		return nil, err
	}
	st, err := loadUndo(p.root)
	if err != nil {
		return nil, err
	}
	entry, ok := st.Entries[p.file]
	if !ok {
		return nil, fmt.Errorf("nothing to undo for %s (the last revision was made by a different run or was already undone)", p.file)
	}

	if err := os.WriteFile(p.file, []byte(entry.Original), 0o644); err != nil {
		return nil, fmt.Errorf("restore %s: %w", p.file, err)
	}
	if err := os.RemoveAll(entry.Archived); err != nil {
		return nil, fmt.Errorf("remove %s: %w", entry.Archived, err)
	}

	// Keep the remaining archives' version lists in sync with the restored page.
	doc, err := fm.Parse(entry.Original)
	if err == nil {
		if err := syncArchives(p.revDir, doc.Revisions()); err != nil {
			return nil, err
		}
	}

	delete(st.Entries, p.file)
	if err := writeUndo(p.root, st); err != nil {
		return nil, err
	}
	return &Result{File: p.file, IsBundle: p.isBundle, Archived: entry.Archived, Entries: doc.Revisions()}, nil
}

func writers(o Options) (io.Writer, io.Writer) {
	out, errw := o.Stdout, o.Stderr
	if out == nil {
		out = os.Stdout
	}
	if errw == nil {
		errw = os.Stderr
	}
	return out, errw
}

func resolve(target string) (*paths, error) {
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", target, err)
	}

	p := &paths{}
	switch {
	case st.IsDir():
		idx := filepath.Join(abs, "index.md")
		if _, err := os.Stat(idx); err != nil {
			return nil, fmt.Errorf("%s is a directory but has no index.md (only page bundles are supported)", abs)
		}
		p.file, p.isBundle = idx, true
		p.dir = abs
		p.name = filepath.Base(abs)
		// The archive dir is a sibling of the bundle: putting it inside would
		// make Hugo read the copies as bundle resources.
		p.revDir = filepath.Join(filepath.Dir(abs), p.name+".revisions")
	default:
		if !strings.EqualFold(filepath.Ext(abs), ".md") {
			return nil, fmt.Errorf("%s: expected a .md file or a page bundle directory", abs)
		}
		p.file, p.isBundle = abs, false
		p.dir = filepath.Dir(abs)
		p.name = strings.TrimSuffix(filepath.Base(abs), ".md")
		p.revDir = filepath.Join(p.dir, p.name+".revisions")
	}

	p.root, err = findRoot(p.dir)
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(filepath.Join(p.root, "content"), p.file)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(rel, "..") {
		return nil, fmt.Errorf("%s is outside %s", p.file, filepath.Join(p.root, "content"))
	}
	p.relPath = filepath.ToSlash(rel)
	return p, nil
}

func findRoot(start string) (string, error) {
	dir := start
	for {
		for _, name := range []string{"hugo.toml", "hugo.yaml", "hugo.yml", "config.toml", "config.yaml", "config.yml"} {
			if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no Hugo project root found (looked for hugo.toml/config.toml upwards from the content file)")
		}
		dir = parent
	}
}

// resolveBase decides the live page URL. Order matters: an explicit flag wins,
// then the URL this tool wrote last time, then Hugo's own permalink, and only
// then a path heuristic (which is wrong for date-based permalinks — hence the
// warning).
func resolveBase(p *paths, entries []fm.Revision, flagURL string, stdout, stderr io.Writer) (string, string, error) {
	if flagURL != "" {
		return withSlash(flagURL), "flag", nil
	}
	if n := len(entries); n > 0 && entries[n-1].URL != "" {
		return withSlash(entries[n-1].URL), "front matter", nil
	}
	if got, err := fromHugoList(p); err == nil && got != "" {
		return withSlash(got), "hugo list all", nil
	} else if err != nil {
		fmt.Fprintf(stderr, "warning: could not ask Hugo for the page URL (%v)\n", err)
	}
	heuristic := "/" + strings.TrimSuffix(strings.TrimPrefix(path.Dir(p.relPath), "."), "/") + "/"
	if !p.isBundle {
		heuristic = "/" + strings.TrimSuffix(path.Dir(p.relPath), ".") + "/" + strings.TrimSuffix(path.Base(p.relPath), ".md") + "/"
	}
	heuristic = strings.ReplaceAll(heuristic, "//", "/")
	fmt.Fprintf(stderr, "warning: assuming the page URL is %q; if your [permalinks] change the path, pass --base-url\n", heuristic)
	return heuristic, "path heuristic", nil
}

// fromHugoList matches this page by its exact content path. Substring matching
// picks the wrong row as soon as a second section contains a same-named page.
func fromHugoList(p *paths) (string, error) {
	cmd := exec.Command("hugo", "list", "all")
	cmd.Dir = p.root
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("hugo list all: %w", err)
	}
	r := csv.NewReader(strings.NewReader(string(out)))
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		return "", fmt.Errorf("parse hugo list all output: %w", err)
	}
	if len(recs) < 2 {
		return "", errors.New("hugo list all returned no pages")
	}
	want, err := indexOf(recs[0], "path", "permalink")
	if err != nil {
		return "", err
	}
	for _, rec := range recs[1:] {
		if len(rec) <= want["path"] || len(rec) <= want["permalink"] {
			continue
		}
		if normalizeListPath(rec[want["path"]]) != path.Clean(p.relPath) {
			continue
		}
		u, err := url.Parse(strings.TrimSpace(rec[want["permalink"]]))
		if err != nil {
			return "", fmt.Errorf("parse permalink %q: %w", rec[want["permalink"]], err)
		}
		// u.Path is already percent-decoded: /2024/01/10/中文标题/
		return u.Path, nil
	}
	return "", fmt.Errorf("page %s not found in hugo list all (is it a draft or excluded?)", p.relPath)
}

// normalizeListPath turns Hugo's `content/posts/foo.md` into `posts/foo.md`,
// which is what paths.relPath holds. Without this every lookup silently misses
// and the tool falls back to a URL guess.
func normalizeListPath(p string) string {
	p = filepath.ToSlash(strings.TrimSpace(p))
	p = strings.TrimPrefix(p, "./")
	p = strings.TrimPrefix(p, "content/")
	return path.Clean(p)
}

func indexOf(header []string, keys ...string) (map[string]int, error) {
	idx := map[string]int{}
	for i, col := range header {
		idx[strings.TrimSpace(col)] = i
	}
	for _, k := range keys {
		if _, ok := idx[k]; !ok {
			return nil, fmt.Errorf("column %q missing from hugo list all output", k)
		}
	}
	return idx, nil
}

func writeArchive(p *paths, src []byte, label, archiveURL string, entries []fm.Revision) (string, error) {
	doc, err := fm.Parse(string(src))
	if err != nil {
		return "", err
	}
	doc.SetSchema(Schema)
	doc.SetCurrentLabel(label)
	doc.SetURL(archiveURL)
	doc.SetBuildNeverList()
	doc.SetRevisions(entries)
	if isDateLabel(label) {
		doc.SetLastmod(label + "T00:00:00Z")
	}

	if !p.isBundle {
		target := filepath.Join(p.revDir, label+".md")
		if _, err := os.Stat(target); err == nil {
			return "", fmt.Errorf("%s already exists", target)
		}
		if err := os.WriteFile(target, doc.Bytes(), 0o644); err != nil {
			return "", err
		}
		return target, nil
	}

	dest := filepath.Join(p.revDir, label)
	if _, err := os.Stat(dest); err == nil {
		return "", fmt.Errorf("%s already exists", dest)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return "", err
	}
	// Copy the bundle's assets (everything but index.md).
	err = filepath.Walk(p.dir, func(srcPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(p.dir, srcPath)
		if err != nil {
			return err
		}
		if rel == "." || rel == "index.md" {
			return nil
		}
		dst := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dest, "index.md"), doc.Bytes(), 0o644); err != nil {
		return "", err
	}
	return dest, nil
}

// syncArchives rewrites the version list inside every archived copy.
func syncArchives(revDir string, entries []fm.Revision) error {
	ents, err := os.ReadDir(revDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range ents {
		target := filepath.Join(revDir, e.Name())
		if e.IsDir() {
			target = filepath.Join(target, "index.md")
		} else if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := os.ReadFile(target)
		if err != nil {
			continue
		}
		doc, err := fm.Parse(string(b))
		if err != nil {
			continue
		}
		doc.SetSchema(Schema)
		doc.SetRevisions(entries)
		if err := os.WriteFile(target, doc.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// mergeEntries records the frozen version at its archive URL, refreshes the
// live version at the base URL, and keeps the list sorted by label.
func mergeEntries(entries []fm.Revision, archivedLabel, archivedURL, liveLabel, base string) []fm.Revision {
	out := make([]fm.Revision, 0, len(entries)+2)
	put := func(label, u string) {
		for i := range out {
			if out[i].Label == label {
				out[i].URL = u
				return
			}
		}
		out = append(out, fm.Revision{Label: label, URL: u})
	}
	for _, e := range entries {
		if e.Label == archivedLabel || e.Label == liveLabel {
			continue
		}
		put(e.Label, e.URL)
	}
	put(archivedLabel, archivedURL)
	put(liveLabel, base)
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out
}

func datePart(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 10 && isDateLabel(v[:10]) {
		return v[:10]
	}
	return ""
}

func isDateLabel(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func withSlash(u string) string {
	if !strings.HasPrefix(u, "/") {
		u = "/" + u
	}
	if !strings.HasSuffix(u, "/") {
		u += "/"
	}
	return u
}

// --- undo state -------------------------------------------------------------

const stateDir = ".hugo-revise"

type undoEntry struct {
	Archived string `json:"archived"`
	Original string `json:"original"`
}

type undoState struct {
	Entries map[string]undoEntry `json:"entries"`
}

func statePath(root string) string { return filepath.Join(root, stateDir, "undo.json") }

func loadUndo(root string) (*undoState, error) {
	st := &undoState{Entries: map[string]undoEntry{}}
	b, err := os.ReadFile(statePath(root))
	if os.IsNotExist(err) {
		return st, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, st); err != nil {
		return nil, fmt.Errorf("read %s: %w", statePath(root), err)
	}
	if st.Entries == nil {
		st.Entries = map[string]undoEntry{}
	}
	return st, nil
}

func writeUndo(root string, st *undoState) error {
	if len(st.Entries) == 0 {
		_ = os.Remove(statePath(root))
		return nil
	}
	if err := os.MkdirAll(filepath.Join(root, stateDir), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(root), b, 0o644)
}

func saveUndo(root, file, archived, original string) error {
	st, err := loadUndo(root)
	if err != nil {
		return err
	}
	st.Entries[file] = undoEntry{Archived: archived, Original: original}
	return writeUndo(root, st)
}
