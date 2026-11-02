package fm

import (
	"strings"
	"testing"
)

const yamlPage = `---
title: "测试文章"
date: 2024-01-10T11:00:00+08:00
tags: ["Hugo", "中文"]
toc: true
---

正文第一段。
`

func TestParseKeepsBodyAndKeys(t *testing.T) {
	d, err := Parse(yamlPage)
	if err != nil {
		t.Fatal(err)
	}
	if d.Format != YAML {
		t.Fatalf("format = %v, want YAML", d.Format)
	}
	if got := d.Date(); got != "2024-01-10T11:00:00+08:00" {
		t.Fatalf("date = %q", got)
	}
	out := string(d.Bytes())
	if !strings.Contains(out, `title: "测试文章"`) {
		t.Fatalf("title line lost:\n%s", out)
	}
	if !strings.Contains(out, "\n正文第一段。\n") {
		t.Fatalf("body lost:\n%s", out)
	}
}

// The bug that broke the previous tool: writing the list twice duplicated the
// key and Hugo refused to build ("mapping key already defined").
func TestSetRevisionsTwiceIsIdempotent(t *testing.T) {
	d, err := Parse(yamlPage)
	if err != nil {
		t.Fatal(err)
	}
	revs := []Revision{{Label: "2024-01-10", URL: "/posts/foo/revisions/2024-01-10/"}}
	for i := 0; i < 3; i++ {
		d.SetSchema(1)
		d.SetRevisions(revs)
	}
	out := string(d.Bytes())
	if n := strings.Count(out, "revisions:"); n != 1 {
		t.Fatalf("revisions block written %d times:\n%s", n, out)
	}
	if n := strings.Count(out, "revisions_schema:"); n != 1 {
		t.Fatalf("revisions_schema written %d times:\n%s", n, out)
	}
	got := d.Revisions()
	if len(got) != 1 || got[0].Label != "2024-01-10" || got[0].URL != "/posts/foo/revisions/2024-01-10/" {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestClearRemovesBlockButKeepsNeighbours(t *testing.T) {
	d, _ := Parse(yamlPage)
	d.SetBuildNeverList()
	d.SetURL("/x/")
	d.Clear("build")
	d.Clear("url")
	out := string(d.Bytes())
	if strings.Contains(out, "build") || strings.Contains(out, "url") {
		t.Fatalf("managed keys survived Clear:\n%s", out)
	}
	if !strings.Contains(out, "toc: true") || !strings.Contains(out, `tags: ["Hugo", "中文"]`) {
		t.Fatalf("neighbour keys lost:\n%s", out)
	}
}

func TestCRLFPreserved(t *testing.T) {
	src := strings.ReplaceAll(yamlPage, "\n", "\r\n")
	d, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	d.SetLastmod("2026-11-02T10:00:00+08:00")
	out := string(d.Bytes())
	if strings.Count(out, "\r\n") != strings.Count(out, "\n") {
		t.Fatalf("mixed line endings after rewrite:\n%q", out)
	}
	if !strings.Contains(out, "lastmod: 2026-11-02T10:00:00+08:00") {
		t.Fatalf("lastmod not written:\n%s", out)
	}
}

func TestTOML(t *testing.T) {
	src := "+++\ntitle = \"测试\"\ndate = 2024-01-10\n+++\n\n正文\n"
	d, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if d.Format != TOML {
		t.Fatalf("format = %v, want TOML", d.Format)
	}
	revs := []Revision{
		{Label: "2024-01-10", URL: "/posts/foo/revisions/2024-01-10/"},
		{Label: "2026-11-02", URL: "/posts/foo/"},
	}
	d.SetSchema(1)
	d.SetRevisions(revs)
	d.SetBuildNeverList()
	out := string(d.Bytes())
	if strings.Count(out, "revisions_schema") != 1 {
		t.Fatalf("schema written twice:\n%s", out)
	}
	got, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parsing TOML output failed: %v\n%s", err, out)
	}
	if r := got.Revisions(); len(r) != 2 || r[1].URL != "/posts/foo/" {
		t.Fatalf("TOML revisions round-trip: %+v", r)
	}
	if !strings.Contains(out, `build = {list = "never", render = true}`) {
		t.Fatalf("build not written:\n%s", out)
	}
}

func TestNoFrontMatter(t *testing.T) {
	if _, err := Parse("# 没有 front matter\n"); err == nil {
		t.Fatal("expected an error for a file without front matter")
	}
}

func TestSetCurrentLabelIsIdempotent(t *testing.T) {
	d, _ := Parse(yamlPage)
	for i := 0; i < 3; i++ {
		d.SetCurrentLabel("2026-11-03")
	}
	out := string(d.Bytes())
	if n := strings.Count(out, "revision_current:"); n != 1 {
		t.Fatalf("revision_current written %d times:\n%s", n, out)
	}
	if got := d.CurrentLabel(); got != "2026-11-03" {
		t.Fatalf("CurrentLabel = %q", got)
	}
}
