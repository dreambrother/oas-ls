# oas-ls

A language server (LSP) for OpenAPI specs, written in Go.

**Status: early development.** Version `0.1.0`, unreleased. APIs and behavior may change.

## Features

- [x] **Go to definition** for `$ref` — line-based (precise cursor targeting not yet implemented)
- [ ] **Rename** a symbol and update all its `$ref`s
- [ ] **Find usages / references**
