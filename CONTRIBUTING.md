# Contributing to GolangAnalyzer

Thanks for helping improve **GolangAnalyzer** (`goproject-all-in-one`).

## Development setup

```bash
git clone https://github.com/Rigsto/goproject-all-in-one.git
cd goproject-all-in-one
go test ./...
go vet ./...
go build -o golanganalyzer .
```

Try the sample project:

```bash
./golanganalyzer flatten --source testdata/sample --out result/flatten.txt
./golanganalyzer minify  --source testdata/sample --out result/minify.txt
```

## Guidelines

- Keep changes **behavior-preserving** unless the issue/PR explicitly changes dump formats or CLI flags.
- Prefer **SOLID / DRY / KISS** — small helpers over new frameworks.
- Add or extend table tests for `minify.Strip` and ignore matching when you touch those paths.
- Run `go test ./...` and `go vet ./...` before opening a PR.
- Do not commit binaries, `result/` dumps, or IDE folders (see `.gitignore`).

## Pull requests

1. Fork and branch from `main`.
2. Make a focused change (one concern per PR when practical).
3. Describe *why* in the PR body; link related issues.
4. Ensure CI is green.

## Reporting bugs

Use the Bug report issue template. Include OS, Go version, the exact command, and a minimal source tree or config if you can.
