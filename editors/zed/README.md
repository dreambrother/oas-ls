# oas-ls — Zed extension

Zed extension that runs the [`oas-ls`](../../README.md) language server for
**YAML** files. It provides OpenAPI `$ref` navigation (go to definition) and is
an *additional* YAML language server: `yaml-language-server` keeps working.

The language server itself lives in the repository root; this directory only
contains the Zed extension (a small Rust program compiled to
`wasm32-wasip2`) that locates and launches the `oas-ls` binary.

## Requirements

- The `oas-ls` binary on `PATH` (see below), or an explicit path in settings.
- Rust via [rustup](https://rustup.rs) — required by Zed to compile the
  extension. Zed installs the `wasm32-wasip2` target automatically.

### Immutable Fedora (Silverblue / Kinoite)

The host has no C toolchain, but Cargo still needs a host linker to build
proc-macro and build-script crates (the wasm artifact itself is linked with the
bundled `rust-lld`). Use a compiler from a toolbox, e.g. `~/.local/bin/toolbox-cc`:

```sh
#!/bin/sh
exec /usr/bin/toolbox run -c <toolbox-with-gcc> gcc "$@"
```

and point Cargo at it in `~/.cargo/config.toml`:

```toml
[target.x86_64-unknown-linux-gnu]
linker = "/home/<user>/.local/bin/toolbox-cc"
```

## Build and install the server

```sh
# from the repository root
go install ./cmd/oas-ls   # binary lands in ~/.go/bin (already on PATH)
# or
go build -o oas-ls ./cmd/oas-ls
```

During development the extension resolves `oas-ls` from `PATH`
(`worktree.which`), so `go install` is enough. Alternatively set
`lsp.oas-ls.binary.path` to a built binary.

## Install the extension (dev)

1. Open the command palette and run `zed: install dev extension`.
2. Select this directory (the one containing `extension.toml`).
3. After editing Rust sources, run `zed: rebuild dev extension`.

Build errors are visible in `zed: open log`; run `zed --foreground` for more
detail.

## Enable the server (opt-in)

`oas-ls` is registered as **opt-in** (`opt_in_languages = ["YAML"]` in
`extension.toml`), so Zed does **not** start it for YAML by default. This is
intentional during early development: enable it explicitly, per project, by
listing it under the language's `language_servers`:

```jsonc
{
  "languages": {
    "YAML": {
      "language_servers": ["oas-ls", "..."]
    }
  }
}
```

Listing `oas-ls` is what turns it on; `"..."` keeps the other YAML servers
(`yaml-language-server`) running. Drop `"..."` if you want *only* `oas-ls`.

Put this in a project's `.zed/settings.json` (or your global settings) only in
repositories where you want the server.

Optional per-server overrides:

```jsonc
{
  "lsp": {
    "oas-ls": {
      "binary": {
        "path": "/home/user/.go/bin/oas-ls",
        "arguments": [],
        "env": {}
      }
    }
  }
}
```

Binary resolution order:

1. `lsp.oas-ls.binary.path` from settings;
2. `worktree.which("oas-ls")` — the `PATH` lookup;
3. otherwise the server fails with an actionable error.

Release auto-download (GitHub releases of `dreambrother/oas-ls` into the
extension work directory) is planned but not implemented yet.

## Notes

- Only `YAML` is registered. JSON specs are not supported by the current
  `oas-ls`.
- The opt-in behavior relies on `opt_in_languages` support in `extension.toml`
  (recent Zed). On an older Zed the field is ignored and the server starts for
  every YAML buffer instead.
- Zed shows the language server's version and memory usage in the status-bar
  LSP popover for locally run servers; nothing needs to be implemented here.
