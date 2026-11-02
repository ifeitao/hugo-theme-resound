// Package fm reads and rewrites Hugo front matter while preserving everything
// it does not own: key order, quoting style and line endings.
//
// The tool never re-serialises front matter through a YAML/TOML marshaller.
// It edits lines, and it removes a managed key together with its block before
// writing a fresh one, so running it twice cannot produce duplicate keys
// (Hugo refuses to build a page with `mapping key already defined`).
package fm

import (
	"errors"
	"fmt"
	"strings"
)

// Format is the front matter syntax of a content file.
type Format int

const (
	Unknown Format = iota
	YAML
	TOML
)

// Revision is one entry of the `revisions` front matter list.
type Revision struct {
	Label string
	URL   string
}

// Doc is a parsed content file: front matter header + raw body.
type Doc struct {
	Format Format
	header []string
	body   string
	eol    string
}

// managed lists the keys this tool owns. They are stripped before a new block
// is written.
var managed = []string{"revisions_schema", "revision_current", "revisions", "url", "build", "lastmod"}

// Parse splits a content file into front matter and body.
func Parse(src string) (*Doc, error) {
	eol := "\n"
	if strings.Contains(src, "\r\n") {
		eol = "\r\n"
	}
	s := strings.TrimPrefix(src, "\ufeff")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	if len(lines) == 0 {
		return nil, errors.New("empty file")
	}

	var delim string
	var format Format
	switch strings.TrimSpace(lines[0]) {
	case "---":
		delim, format = "---", YAML
	case "+++":
		delim, format = "+++", TOML
	default:
		return nil, errors.New("no front matter found (expected --- or +++ on the first line)")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == delim {
			end = i
			break
		}
	}
	if end < 0 {
		return nil, fmt.Errorf("front matter opened with %q but never closed", delim)
	}

	return &Doc{
		Format: format,
		header: append([]string(nil), lines[1:end]...),
		body:   strings.Join(lines[end+1:], "\n"),
		eol:    eol,
	}, nil
}

// Bytes renders the document back to a complete file.
func (d *Doc) Bytes() []byte {
	delim := "---"
	if d.Format == TOML {
		delim = "+++"
	}
	var b strings.Builder
	b.WriteString(delim)
	b.WriteString(d.eol)
	if len(d.header) > 0 {
		b.WriteString(strings.Join(d.header, d.eol))
		b.WriteString(d.eol)
	}
	b.WriteString(delim)
	b.WriteString(d.eol)
	b.WriteString(strings.ReplaceAll(d.body, "\n", d.eol))
	return []byte(b.String())
}

// Date returns the date scalar exactly as written (may be quoted).
func (d *Doc) Date() string {
	v, _ := d.scalar("date")
	return v
}

// HasSchema reports whether this page was written by a tool version we know.
func (d *Doc) HasSchema() bool {
	v, _ := d.scalar("revisions_schema")
	return strings.TrimSpace(v) == "1"
}

// Revisions returns the parsed `revisions` list, in the order written.
func (d *Doc) Revisions() []Revision {
	block := d.keyBlock("revisions")
	if len(block) == 0 {
		return nil
	}
	if d.Format == TOML {
		return parseTOMLRevisions(block)
	}
	return parseYAMLRevisions(block)
}

// SetSchema stamps the contract version.
func (d *Doc) SetSchema(v int) {
	d.Clear("revisions_schema")
	d.header = append(d.header, fmt.Sprintf("revisions_schema: %d", v))
	if d.Format == TOML {
		d.header[len(d.header)-1] = fmt.Sprintf("revisions_schema = %d", v)
	}
}

// SetRevisions replaces the list with the canonical block form.
func (d *Doc) SetRevisions(revs []Revision) {
	d.Clear("revisions")
	if len(revs) == 0 {
		return
	}
	if d.Format == TOML {
		var b strings.Builder
		b.WriteString("revisions = [")
		for i, r := range revs {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "{label = %q, url = %q}", r.Label, r.URL)
		}
		b.WriteString("]")
		d.header = append(d.header, b.String())
		return
	}
	d.header = append(d.header, "revisions:")
	for _, r := range revs {
		d.header = append(d.header,
			"  - label: "+quoteYAML(r.Label),
			"    url: "+quoteYAML(r.URL))
	}
}

// SetURL pins the published path of an archived page.
func (d *Doc) SetURL(u string) {
	d.Clear("url")
	if d.Format == TOML {
		d.header = append(d.header, fmt.Sprintf("url = %q", u))
		return
	}
	d.header = append(d.header, "url: "+quoteYAML(u))
}

// SetLastmod stamps the modification time (unquoted RFC 3339, which Hugo parses).
func (d *Doc) SetLastmod(v string) {
	d.Clear("lastmod")
	if d.Format == TOML {
		d.header = append(d.header, "lastmod = "+v)
		return
	}
	d.header = append(d.header, "lastmod: "+v)
}

// SetCurrentLabel records which version this file is. The theme matches the
// selector on this value instead of comparing URLs, because RelPermalink is
// percent-encoded for non-ASCII paths and front matter URLs are not.
func (d *Doc) SetCurrentLabel(label string) {
	d.Clear("revision_current")
	if d.Format == TOML {
		d.header = append(d.header, fmt.Sprintf("revision_current = %q", label))
		return
	}
	d.header = append(d.header, "revision_current: "+quoteYAML(label))
}

// CurrentLabel returns the label this file claims to be.
func (d *Doc) CurrentLabel() string {
	v, _ := d.scalar("revision_current")
	return v
}

// SetBuildNeverList publishes the page but keeps it out of every list,
// taxonomy, RSS feed and search index.
func (d *Doc) SetBuildNeverList() {
	d.Clear("build")
	if d.Format == TOML {
		d.header = append(d.header, `build = {list = "never", render = true}`)
		return
	}
	d.header = append(d.header, "build:", "  list: never", "  render: true")
}

// Clear removes a managed key together with its indented block.
func (d *Doc) Clear(key string) {
	out := make([]string, 0, len(d.header))
	for i := 0; i < len(d.header); i++ {
		if !d.isKey(d.header[i], key) {
			out = append(out, d.header[i])
			continue
		}
		// Skip the key line, then everything that belongs to its value.
		i = d.blockEnd(i)
	}
	d.header = out
}

// blockEnd returns the index of the last line belonging to the key at i.
func (d *Doc) blockEnd(i int) int {
	last := i
	for j := i + 1; j < len(d.header); j++ {
		line := d.header[j]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			// A blank line ends a scalar, but not an open multi-line array.
			if d.Format == TOML && strings.Contains(d.header[last], "[") && !strings.Contains(d.header[last], "]") {
				last = j
				continue
			}
			break
		}
		indented := strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
		listItem := strings.HasPrefix(trimmed, "- ")
		// TOML multi-line arrays continue until the closing bracket.
		if d.Format == TOML && strings.Contains(d.header[last], "[") && !strings.Contains(d.header[last], "]") {
			last = j
			continue
		}
		if !indented && !listItem {
			break
		}
		last = j
	}
	return last
}

// keyBlock returns the lines of a key (key line first, continuation lines after).
func (d *Doc) keyBlock(key string) []string {
	for i := 0; i < len(d.header); i++ {
		if d.isKey(d.header[i], key) {
			return d.header[i : d.blockEnd(i)+1]
		}
	}
	return nil
}

// isKey reports whether line assigns key at column 0.
func (d *Doc) isKey(line, key string) bool {
	if d.Format == TOML {
		return strings.HasPrefix(line, key+" ") || strings.HasPrefix(line, key+"=")
	}
	return strings.HasPrefix(line, key+":")
}

// scalar returns the value of a single-line key, unquoted.
func (d *Doc) scalar(key string) (string, bool) {
	for _, line := range d.header {
		if !d.isKey(line, key) {
			continue
		}
		var v string
		if d.Format == TOML {
			v = strings.TrimPrefix(line, key)
			v = strings.TrimSpace(strings.TrimPrefix(v, "="))
		} else {
			v = strings.TrimSpace(strings.TrimPrefix(line, key+":"))
		}
		return unquote(v), true
	}
	return "", false
}

// quoteYAML always quotes a scalar. Labels such as 2024-01-10 are parsed as
// timestamps when left bare, and Hugo hands the template a time.Time instead of
// a string — which silently breaks every label comparison.
func quoteYAML(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func valueOf(entry, key string) string {
	// Matches both `label: 2024` (YAML) and `label = "2024"` (TOML).
	if idx := strings.Index(entry, key+":"); idx >= 0 {
		return unquote(entry[idx+len(key)+1:])
	}
	if idx := strings.Index(entry, key+" ="); idx >= 0 {
		return unquote(entry[idx+len(key)+2:])
	}
	if idx := strings.Index(entry, key+"="); idx >= 0 {
		return unquote(entry[idx+len(key)+1:])
	}
	return ""
}

func parseYAMLRevisions(block []string) []Revision {
	var out []Revision
	var cur *Revision
	for _, line := range block[1:] {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "- ") {
			if cur != nil {
				out = append(out, *cur)
			}
			cur = &Revision{}
			t = strings.TrimSpace(strings.TrimPrefix(t, "- "))
		}
		if cur == nil {
			continue
		}
		if v := valueOf(t, "label"); v != "" {
			cur.Label = v
		}
		if v := valueOf(t, "url"); v != "" {
			cur.URL = v
		}
	}
	if cur != nil {
		out = append(out, *cur)
	}
	return out
}

func parseTOMLRevisions(block []string) []Revision {
	var out []Revision
	joined := strings.Join(block, " ")
	// Inline tables: {label = "x", url = "y"}
	for _, part := range strings.Split(joined, "}") {
		open := strings.Index(part, "{")
		if open < 0 {
			continue
		}
		fields := part[open+1:]
		label := valueOf(fields, "label")
		url := valueOf(fields, "url")
		if label == "" && url == "" {
			continue
		}
		out = append(out, Revision{Label: label, URL: url})
	}
	return out
}
