# How to use GolangAnalyzer

End-to-end flows for building a project dump and reviewing it with an AI. Skim the [README](README.md) for what the tool does; this guide is for running it against *your* needs.

## 1. Build the CLI

From the GolangAnalyzer repo root:

```bash
go build -o golanganalyzer .
```

On Windows PowerShell you can use:

```powershell
go build -o golanganalyzer.exe .
```

Optional: put the binary on your `PATH`, or call it with `./golanganalyzer` / `.\golanganalyzer.exe` from this directory.

Verify:

```bash
./golanganalyzer --help
./golanganalyzer flatten --help
./golanganalyzer minify --help
```

## 2. Pick flatten vs minify

| You need… | Use |
|-----------|-----|
| Exact source, comments kept, smallest chance of citation mistakes | **`flatten`** |
| Smaller dump for large repos / tight token limits | **`minify`** |
| Trying the tool on the included sample | Either — see [Try the sample](#3-try-the-sample) |

Rule of thumb: start with **minify** for big trees; switch to **flatten** for a package or feature you want to inspect line-by-line with comments intact.

## 3. Try the sample

```bash
./golanganalyzer flatten --source testdata/sample --out result/sample-flatten.txt
./golanganalyzer minify  --source testdata/sample --out result/sample-minify.txt
```

Open the files under `result/` and compare formats before pointing `--source` at a real project.

## 4. Analyze your own project

### Minimal run

```bash
./golanganalyzer flatten --source /path/to/your/project
# → result/flatten.txt

./golanganalyzer minify --source /path/to/your/project
# → result/minify.txt
```

`--source` must be a **directory**. Paths inside the dump are relative to that root.

### Custom output path

```bash
./golanganalyzer minify \
  --source ~/work/my-api \
  --out ~/reviews/my-api-minify.txt
```

Parent directories for `--out` are created if missing.

### Custom config

Copy and edit `config.yaml`, then pass it:

```bash
./golanganalyzer --config my-config.yaml minify --source ~/work/my-api
```

Useful edits:

- **Narrow extensions** — e.g. only `go` for a Go service, or add `ts` / `tsx` for a frontend.
- **Ignore noise** — generated dirs (`vendor/`, `node_modules/`, `dist/`), build output, fixtures you do not want reviewed.
- **Exclude secrets** — keep `.env`, key material, and credential paths out of `extensions` / behind `ignore` so they never enter the dump.

Example focused on Go only:

```yaml
extensions:
  - go
ignore:
  - .git/
  - vendor/
  - bin/
  - result/
  - "*_test.go"   # optional: skip tests if you only want production code
```

> Note: ignore uses simple path/glob matching (directory names, path prefixes, basename globs). Prefer listing real directories and globs you know exist in the tree.

## 5. Review with an AI (recommended flow)

1. **Generate the dump** (`flatten` or `minify`) as above.
2. **Choose a prompt** from the folder that matches your dump:
   - Flatten dump → `prompt/flatten/`
   - Minify dump → `prompt/minify/`
   Topics in each folder:
   - `analysis.txt` — clean code / maintainability
   - `security.txt` — vulnerabilities and risky patterns
   - `performance.txt` — hotspots and inefficient patterns
   - `error-handling.txt` — swallowed errors, panics, partial failure
   - `concurrency.txt` — races, goroutine leaks, lock/channel issues
   - `explain.txt` — onboarding tour of the dump
3. **In your AI chat**: paste the full prompt first, then paste the dump (or attach the dump file if your tool supports files).
4. **Read findings** — citations should use **original** file paths and line numbers (minify prompts tell the model to apply `@MAP`).

You can run the same dump through more than one topic (security pass, then performance pass) without regenerating — keep using the matching `flatten/` or `minify/` folder.

### Checklist before you paste

- [ ] Dump does not include secrets, private keys, or production configs with credentials
- [ ] Extensions/ignores match what you want reviewed
- [ ] File size fits your model’s context (if not, tighten `ignore` or use `minify` / a subdirectory as `--source`)

## 6. Flows by need

### A. Quick clean-code pass on a feature folder

```bash
./golanganalyzer flatten \
  --source ~/work/my-api/internal/billing \
  --out result/billing-flatten.txt
```

Use `prompt/flatten/analysis.txt` + that dump. Pointing `--source` at a subdirectory keeps the dump small and focused.

### B. Security review of a whole service (token-conscious)

```bash
# edit config.yaml extensions/ignore first if needed
./golanganalyzer minify \
  --source ~/work/my-api \
  --out result/my-api-minify.txt
```

Use `prompt/minify/security.txt`. Prefer minify so more of the tree fits; rely on `@MAP` for original lines.

### C. Performance pass on hot packages only

```bash
./golanganalyzer minify \
  --source ~/work/my-api/internal/handler \
  --out result/handler-perf.txt
```

Use `prompt/minify/performance.txt`. Repeat for other hot packages (`repo/`, `worker/`, etc.) if one dump is still too large.

### D. Multi-language monorepo

Extend `extensions` (e.g. `go`, `py`, `js`, `ts`, `yaml`) and ignore language-specific junk:

```yaml
ignore:
  - .git/
  - node_modules/
  - vendor/
  - __pycache__/
  - .venv/
  - dist/
  - build/
```

Then minify the monorepo root, or flatten one language subtree at a time.

### E. Re-run after code changes

Regenerate the same `--out` path after edits, then paste the new dump with the same prompt. No extra steps — collection always walks the live tree.

## 7. Reading minify maps (optional)

If you want to verify a citation yourself:

```text
@FILE internal/api/server.go
@ORIG 1-120
@MAP 1>5,2>6,3>10
line A   ← dump line 1 → original line 5
line B   ← dump line 2 → original line 6
line C   ← dump line 3 → original line 10
```

Flatten dumps need no map: dump line numbers match the source.

## 8. Common issues

| Symptom | What to do |
|---------|------------|
| `--source is required` | Pass `--source <dir>` |
| `source is not a directory` | Point at a folder, not a single file |
| `config: extensions must not be empty` | Add at least one extension in the YAML |
| Empty or tiny dump | Check `extensions` and `ignore`; confirm files exist under `--source` |
| Dump too large for the model | Use `minify`, narrow `--source`, or expand `ignore` |
| Wrong line in the AI reply | For minify, use `prompt/minify/*.txt` so `@MAP` is applied; for flatten use `prompt/flatten/*.txt` |

## 9. One-liner cheat sheet

```bash
# Sample
./golanganalyzer flatten --source testdata/sample

# Your project — full source
./golanganalyzer flatten --source /path/to/project --out result/flatten.txt

# Your project — compressed + maps
./golanganalyzer minify  --source /path/to/project --out result/minify.txt

# Custom rules
./golanganalyzer --config my-config.yaml minify --source /path/to/project
```

Then: open `prompt/<flatten|minify>/<topic>.txt` (match your dump) → paste prompt → paste dump → review.
