# tools/revise

Version history for Hugo content. Companion to the Resound theme's version
selector — **optional**: without it the theme builds and renders exactly as
before, just without a selector.

It runs at *writing* time. It is never part of `hugo`, Netlify, or CI builds, and
needs no secrets.

## Install

```bash
# with Go 1.22+
go install github.com/ifeitao/hugo-theme-resound/tools/cmd/revise@latest

# or download the binary for your platform from the repository's Releases,
# put it on your PATH and call it `revise` (revise.exe on Windows)
```

Run from anywhere inside the Hugo project; the tool walks up to find the site
root (`hugo.toml`, `config.toml`, …).

## Usage

```bash
revise [--base-url URL] [--label LABEL] <content/posts/foo.md | content/posts/foo-bundle>
revise undo <content/posts/foo.md | content/posts/foo-bundle>
```

| Flag | Meaning |
|---|---|
| `--label` | Label for the new version. Default: today (`2006-01-02`). Use e.g. `2026-11-03-2` for a second revision on the same day. Labels sort chronologically as strings. |
| `--base-url` | The page's current URL, `/posts/foo/`. Needed when `[permalinks]` rewrites the path and Hugo cannot be asked. |

Commit before revising:

```bash
git add -A && git commit -m "before revision"
revise content/posts/foo.md
git add -A && git commit -m "archive the previous version"
```

`revise undo <path>` restores the page byte for byte and deletes the archive
created by the last run on that page. Its state lives in `.hugo-revise/undo.json`
at the project root — add `.hugo-revise/` to the site's `.gitignore`.

## What it writes

```
content/posts/
├── foo.md                        # live page
└── foo.revisions/
    ├── 2024-01-10.md             # archived copy, published at <base>revisions/2024-01-10/
    └── 2026-11-03.md
```

Page bundles get a directory per version, with every asset copied except
`index.md`. The archive directory is a **sibling** of the page, never inside the
bundle: a bundle root may contain only one `.md`.

Front matter contract (schema 1):

```yaml
revisions_schema: 1
revision_current: "2026-11-03"          # which version this file is
revisions:                              # label -> URL
  - label: "2024-01-10"
    url: "/2024/01/10/foo/revisions/2024-01-10/"
  - label: "2026-11-03"
    url: "/2024/01/10/foo/"
```

Archived pages additionally get `url` (pinned path) and
`build: {list: never, render: true}` so they are reachable but absent from
lists, taxonomies, RSS, search and the sitemap.

## Design decisions worth knowing

- **`date` is never rewritten.** Only `lastmod` changes. Under date-based
  permalinks (`posts = '/:year/:month/:day/:title/'`) rewriting `date` moves the
  live URL and breaks every inbound link; the previous tool did exactly that.
- **URLs are stored, not derived.** The theme iterates `revisions`; it never
  splits `RelPermalink` on `/revisions/`. Two independent derivations of the same
  base URL always drift eventually.
- **`selected` is matched on labels**, because `RelPermalink` is percent-encoded
  for non-ASCII paths while front matter URLs are decoded.
- **`hugo list all` permalinks are percent-decoded before use.** Writing the
  encoded form into `url:` publishes a literally `%E9%9D%99…` directory that no
  server can serve.
- **`hugo list all` rows are matched on the exact content path.** Substring
  matching picks `en/posts/x.md` when you asked for `posts/x.md`.
- **Front matter is edited line by line**, not re-serialised: key order, quoting
  and line endings survive. Managed keys are stripped before being rewritten, so
  a second run cannot produce `mapping key already defined`.
- **Labels are quoted.** A bare `2024-01-10` is parsed as a timestamp, and Hugo
  hands the template a `time.Time` instead of a string.

## Tests

```bash
go test ./tools/...     # from the repository root
```

`internal/fm` covers line-ending preservation, idempotent writes, TOML and
duplicate-key regressions. `internal/archive` covers the two-revision + undo
cycle and bundle asset copying.

## Releases

`.github/workflows/release.yml` cross-compiles linux/amd64, linux/arm64,
darwin/amd64, darwin/arm64 and windows/amd64 on every `v*` tag. The binary is
pure Go (`CGO_ENABLED=0`), so each artifact is a single self-contained file.
