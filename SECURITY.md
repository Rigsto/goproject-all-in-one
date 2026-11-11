# Security Policy

## Supported versions

Security fixes are applied on the latest `main` branch and the newest GitHub Release when one exists.

## Reporting a vulnerability

Please **do not** open a public issue for security-sensitive reports.

Email or contact the maintainer via GitHub ([@Rigsto](https://github.com/Rigsto)) with:

- A short description of the issue
- Steps to reproduce
- Impact (e.g. unexpected file read outside `--source`, crash on crafted input)

You should receive an acknowledgment when possible. Please give reasonable time for a fix before public disclosure.

## Scope notes

GolangAnalyzer reads local files under `--source` according to `config.yaml`. Treat untrusted trees like any other local tooling: run it only on code you trust, or in an isolated environment.
