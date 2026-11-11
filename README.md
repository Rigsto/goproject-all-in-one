# GolangAnalyzer

CLI that packs a source tree into one text dump so you can paste it into an AI chat for code review.

It walks a project, keeps only the extensions you care about, skips ignore rules, then writes either a **full-source flatten** or a **comment-stripped minify** (with line maps back to the original files). Ready-made prompts under `prompt/flatten/` and `prompt/minify/` match each dump format and cover clean-code, security, performance, error handling, concurrency, and codebase explainers.

## Why this exists

AI reviewers work better when they see the whole project in one place, with paths and line numbers they can cite. GolangAnalyzer builds that dump locally — you choose fidelity vs. token size, then pair the output with a prompt.

| Mode | Best for | Output |
|------|----------|--------|
| `flatten` | Small/medium trees, exact line citations | Full file bodies with `===== FILE path \| LINES 1-n` headers |
| `minify` | Larger trees / tighter context windows | Comments & blank lines removed; `@MAP` maps dump lines → original lines |

## Requirements

- [Go](https://go.dev/dl/) 1.22+ (module targets a recent toolchain)
- A project directory you want to analyze (any mix of supported extensions)

## Quick start

```bash
# from this repo
go build -o golanganalyzer .

# pack a project (full source)
./golanganalyzer flatten --source /path/to/your/project

# or pack a compressed dump
./golanganalyzer minify --source /path/to/your/project
```

Defaults write to `result/flatten.txt` or `result/minify.txt`. Then open a matching prompt from `prompt/flatten/` or `prompt/minify/`, paste the dump after it, and run the review in your AI tool of choice.

For step-by-step flows (config tweaks, custom output paths, which prompt to pick), see **[HOW-TO-USE.md](HOW-TO-USE.md)**.

## Commands

```text
golanganalyzer [--config path] flatten --source <dir> [--out <file>]
golanganalyzer [--config path] minify  --source <dir> [--out <file>]
```

| Flag | Default | Meaning |
|------|---------|---------|
| `--config` | `config.yaml` | Extensions + ignore rules |
| `--source` | *(required)* | Root directory to walk |
| `--out` | `result/flatten.txt` or `result/minify.txt` | Dump destination |

## Configuration

`config.yaml` controls what gets included:

```yaml
extensions:
  - go
  - js
  - py
  - java
  - yaml
  - yml
ignore:
  - .git/
  - vendor/
  - node_modules/
  - "*.min.js"
  - dist/
  - bin/
  - result/
```

- **extensions** — file suffixes to keep (dots optional; matched case-insensitively).
- **ignore** — gitignore-style paths/globs (directories ending in `/`, name globs like `*.min.js`).

Binary files (NUL byte in content) are skipped automatically.

## Output formats

### Flatten

```text
===== FILE path/to/file.go | LINES 1-42
package main
...
```

Body line *N* is original line *N*. No map needed.

### Minify

```text
@FILE path/to/file.go
@ORIG 1-42
@MAP 1>3,2>4,3>7,...
<stripped lines>
```

`@MAP out>orig` means dump line `out` came from original line `orig`. Review prompts instruct the model to cite **original** lines using this map.

## Prompts

Pick the folder that matches your dump command, then a topic:

| Topic | Flatten | Minify | Focus |
|-------|---------|--------|--------|
| analysis | [`prompt/flatten/analysis.txt`](prompt/flatten/analysis.txt) | [`prompt/minify/analysis.txt`](prompt/minify/analysis.txt) | Clean code, dead code, naming, structure |
| security | [`prompt/flatten/security.txt`](prompt/flatten/security.txt) | [`prompt/minify/security.txt`](prompt/minify/security.txt) | Vulnerabilities, auth gaps, unsafe patterns |
| performance | [`prompt/flatten/performance.txt`](prompt/flatten/performance.txt) | [`prompt/minify/performance.txt`](prompt/minify/performance.txt) | Allocations, hot paths, inefficient loops |
| error-handling | [`prompt/flatten/error-handling.txt`](prompt/flatten/error-handling.txt) | [`prompt/minify/error-handling.txt`](prompt/minify/error-handling.txt) | Swallowed errors, panics, partial failure |
| concurrency | [`prompt/flatten/concurrency.txt`](prompt/flatten/concurrency.txt) | [`prompt/minify/concurrency.txt`](prompt/minify/concurrency.txt) | Races, goroutine leaks, channel/lock issues |
| explain | [`prompt/flatten/explain.txt`](prompt/flatten/explain.txt) | [`prompt/minify/explain.txt`](prompt/minify/explain.txt) | Onboarding tour of the dump |

Each prompt describes **one** input format only — use `flatten/` with a flatten dump, `minify/` with a minify dump.

## Layout

```text
golanganalyzer/
├── main.go              # CLI (cobra): flatten / minify
├── config.yaml          # default include/ignore rules
├── config/              # YAML load + validate
├── flatten/             # full-source dump writer
├── minify/              # strip comments/blanks + line maps
├── internal/collect/    # walk tree, apply extensions & ignores
├── prompt/              # AI review prompts (flatten/ and minify/)
├── testdata/sample/     # tiny sample project
├── HOW-TO-USE.md        # end-to-end usage flows
└── result/              # default output dir (gitignored)
```

## License

Use and modify as needed for your own analysis workflows.
