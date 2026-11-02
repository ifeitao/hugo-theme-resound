package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ifeitao/hugo-theme-resound/tools/internal/fm"
)

func TestMergeEntries(t *testing.T) {
	got := mergeEntries(
		[]fm.Revision{{Label: "2024-01-10", URL: "/p/"}},
		"2024-01-10", "/p/revisions/2024-01-10/",
		"2026-11-02", "/p/",
	)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2: %+v", len(got), got)
	}
	if got[0].URL != "/p/revisions/2024-01-10/" || got[1].URL != "/p/" {
		t.Fatalf("wrong urls: %+v", got)
	}
	if got[0].Label > got[1].Label {
		t.Fatalf("not sorted: %+v", got)
	}
}

func TestDatePart(t *testing.T) {
	cases := map[string]string{
		"2024-01-10T11:00:00+08:00": "2024-01-10",
		"2024-01-10":                "2024-01-10",
		"not a date":                "",
		"":                          "",
	}
	for in, want := range cases {
		if got := datePart(in); got != want {
			t.Fatalf("datePart(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestWithSlash(t *testing.T) {
	if got := withSlash("posts/foo"); got != "/posts/foo/" {
		t.Fatalf("withSlash = %q", got)
	}
}

func TestNormalizeListPath(t *testing.T) {
	cases := map[string]string{
		"content/posts/foo.md":  "posts/foo.md",
		"content/en/posts/x.md": "en/posts/x.md",
		"./content/posts/a.md":  "posts/a.md",
		"posts/a.md":            "posts/a.md",
	}
	for in, want := range cases {
		if got := normalizeListPath(in); got != want {
			t.Fatalf("normalizeListPath(%q) = %q, want %q", in, got, want)
		}
	}
}

// writeSite lays out a minimal Hugo project with one page.
func writeSite(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(root, "hugo.toml"), []byte("baseURL = \"https://e.com/\"\n"), 0o644))
	must(os.MkdirAll(filepath.Join(root, "content", "posts"), 0o755))
	must(os.WriteFile(filepath.Join(root, "content", "posts", "foo.md"), []byte(body), 0o644))
	return root
}

const page = `---
title: "foo"
date: 2024-01-10T09:00:00+08:00
tags: [a, b]
---

正文
`

func TestReviseTwiceThenUndo(t *testing.T) {
	root := writeSite(t, page)
	target := filepath.Join(root, "content", "posts", "foo.md")

	// Explicit labels keep the test independent of the wall clock.
	res, err := Revise(target, Options{BaseURL: "/posts/foo/", Label: "2024-06-01"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Label != "2024-01-10" {
		t.Fatalf("first archive label = %q, want the publish date", res.Label)
	}
	arch := filepath.Join(root, "content", "posts", "foo.revisions", "2024-01-10.md")
	if _, err := os.Stat(arch); err != nil {
		t.Fatalf("archive missing: %v", err)
	}

	live := read(t, target)
	if !strings.Contains(live, "date: 2024-01-10T09:00:00+08:00") {
		t.Fatalf("date must not be rewritten:\n%s", live)
	}
	if strings.Count(live, "revisions:") != 1 {
		t.Fatalf("revisions block duplicated:\n%s", live)
	}

	// Second revision must not duplicate keys and must extend the list.
	afterFirst := read(t, target)
	res2, err := Revise(target, Options{Label: "2024-12-01"})
	if err != nil {
		t.Fatal(err)
	}
	if res2.BaseSource != "front matter" {
		t.Fatalf("second run should reuse the URL from front matter, got %q", res2.BaseSource)
	}
	if len(res2.Entries) != 3 {
		t.Fatalf("entries = %d, want 3: %+v", len(res2.Entries), res2.Entries)
	}
	arch2 := filepath.Join(root, "content", "posts", "foo.revisions", "2024-06-01.md")
	if _, err := os.Stat(arch2); err != nil {
		t.Fatalf("second archive missing (%s): %v", arch2, err)
	}
	first := read(t, arch)
	if strings.Count(first, "revisions:") != 1 || !strings.Contains(first, "2024-12-01") {
		t.Fatalf("older archive not re-synced:\n%s", first)
	}
	if !strings.Contains(first, `url: "/posts/foo/revisions/2024-01-10/"`) {
		t.Fatalf("archive url missing:\n%s", first)
	}
	if !strings.Contains(first, "list: never") {
		t.Fatalf("build options missing:\n%s", first)
	}

	// Undo removes the newest archive and restores the page to its state before
	// the last revision, byte for byte.
	if _, err := Undo(target); err != nil {
		t.Fatal(err)
	}
	if got := read(t, target); got != afterFirst {
		t.Fatalf("page not restored:\n--- got ---\n%s\n--- want ---\n%s", got, afterFirst)
	}
	if _, err := os.Stat(arch2); !os.IsNotExist(err) {
		t.Fatalf("second archive should be gone, err=%v", err)
	}
	if !strings.HasSuffix(read(t, arch), "正文\n") {
		t.Fatalf("first archive damaged by undo")
	}
	if strings.Contains(read(t, arch), "2024-12-01") {
		t.Fatalf("undo left the removed label in the older archive")
	}
}

func TestBundleCopiesAssets(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "content", "posts", "bundle")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, "hugo.toml"), []byte("baseURL = \"https://e.com/\"\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "index.md"), []byte(page), 0o644)
	os.WriteFile(filepath.Join(dir, "pic.png"), []byte("PNG"), 0o644)
	os.MkdirAll(filepath.Join(dir, "nested"), 0o755)
	os.WriteFile(filepath.Join(dir, "nested", "x.txt"), []byte("X"), 0o644)

	if _, err := Revise(dir, Options{BaseURL: "/posts/bundle/"}); err != nil {
		t.Fatal(err)
	}
	archDir := filepath.Join(root, "content", "posts", "bundle.revisions", "2024-01-10")
	for _, want := range []string{"index.md", "pic.png", filepath.Join("nested", "x.txt")} {
		if _, err := os.Stat(filepath.Join(archDir, want)); err != nil {
			t.Fatalf("asset %s not copied: %v", want, err)
		}
	}
	// The bundle itself must stay clean: no archive dir inside it.
	if _, err := os.Stat(filepath.Join(dir, "revisions")); err == nil {
		t.Fatal("archive directory must not live inside the page bundle")
	}
}

func read(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
