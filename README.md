# GolangAnalyzer

[![CI](https://github.com/Rigsto/goproject-all-in-one/actions/workflows/ci.yml/badge.svg)](https://github.com/Rigsto/goproject-all-in-one/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)

**Repo:** [`goproject-all-in-one`](https://github.com/Rigsto/goproject-all-in-one) · **Tool:** `golanganalyzer`

Pack a whole project into **one text dump** for AI code review — with paths and line numbers the model can cite.

```text
your repo  →  golanganalyzer flatten|minify  →  paste into ChatGPT / Claude / local LLM
                 + matching prompt from prompt/
```

| Mode | Best for | Output |
|------|----------|--------|
| `flatten` | Exact citations, comments kept | Full bodies · `===== FILE path \| LINES 1-n` |
| `minify` | Large trees / tight token budgets | Comments & blanks removed · `@MAP` dump→original lines |

Ready-made prompts cover clean-code, security, performance, error handling, concurrency, and codebase explainers.

## Install

**Go install** (needs Go 1.22+):

```bash
go install github.com/Rigsto/goproject-all-in-one@latest
```

The binary is named after the module folder (`goproject-all-in-one`). Rename it if you prefer:

```bash
# Unix
mv "$(go env GOPATH)/bin/goproject-all-in-one" "$(go env GOPATH)/bin/golanganalyzer"

# Windows PowerShell
Move-Item "$(go env GOPATH)\bin\goproject-all-in-one.exe" "$(go env GOPATH)\bin\golanganalyzer.exe"
```

**From source:**

```bash
git clone https://github.com/Rigsto/goproject-all-in-one.git
cd goproject-all-in-one
go build -o golanganalyzer .
```

**Releases:** download a prebuilt binary from [Releases](https://github.com/Rigsto/goproject-all-in-one/releases) once published.

## 30-second demo

```bash
go build -o golanganalyzer .
./golanganalyzer flatten --source testdata/sample --out result/demo-flatten.txt
./golanganalyzer minify  --source testdata/sample --out result/demo-minify.txt
```

**Flatten** (body line *N* = original line *N*):

```text
===== FILE main.go | LINES 1-13
package main
...
func main() {
	// entry
	greet("world")
}
```

**Minify** (comments gone; `@MAP` keeps original line citations):

```text
@FILE main.go
@ORIG 1-13
@MAP 1>1,2>3,3>6,4>7,5>8,6>10,7>12,8>13
package main
import "fmt"
func greet(name string) {
	fmt.Println("hello", name)
}
func main() {
	greet("world")
}
```

## Workflow

1. Build or install the CLI (`golanganalyzer`).
2. Run `flatten` or `minify` on your project (`--source`).
3. Open a prompt from `prompt/flatten/` or `prompt/minify/` (same mode as the dump).
4. Paste **prompt + dump** into your AI tool and ask for the review.

Step-by-step flows (config, custom `--out`, which prompt): **[HOW-TO-USE.md](HOW-TO-USE.md)**.

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
goproject-all-in-one/          # GitHub repository
├── main.go                    # CLI binary: golanganalyzer
├── config.yaml
├── config/ · flatten/ · minify/
├── internal/collect/ · internal/text/
├── prompt/                    # AI review prompts
├── testdata/sample/
├── HOW-TO-USE.md
└── result/                    # default output (gitignored)
```

## Architecture

```text
CLI (main.go)
  │  --config → config.Load
  │  --source / --out
  ├─ flatten → flatten.FlattenWithConfig
  └─ minify  → minify.MinifyWithConfig
                    │
                    ▼
            collect.Collect
              • WalkDir from --source
              • keep cfg.Extensions
              • skip cfg.Ignore (+ binary / NUL)
                    │
         ┌──────────┴──────────┐
         ▼                     ▼
   flatten writer         minify writer
   ===== FILE …           @FILE / @ORIG / @MAP
   full body              Strip + FormatMap
```

| Package | Role |
|---------|------|
| `main` | Cobra CLI; loads config, dispatches flatten/minify |
| [`config`](config/config.go) | Parse/validate `config.yaml` |
| [`internal/collect`](internal/collect/collect.go) | Tree walk, ignore matching, sorted files |
| [`flatten`](flatten/flatten.go) | Full-source dump |
| [`minify`](minify/minify.go) | Stripped dump + [`Strip`](minify/strip.go) / `@MAP` |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and ideas welcome via [Issues](https://github.com/Rigsto/goproject-all-in-one/issues).

## License

[MIT](LICENSE) © Rigsto
