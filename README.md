# oas-ls

A language server (LSP) for OpenAPI specs, written in Go.

**Status: early development.** Version `0.1.0`, unreleased. APIs and behavior may change.

## Supported formats

Only **YAML** specs are supported so far (`.yaml`, `.yml`).

**JSON specs are not supported yet.** The server parses documents with a YAML
parser and detects `$ref` in a line-based way, which does not handle JSON
reliably. Associate only YAML files with `oas-ls` in your editor.

## Features

- [x] **Go to definition** for `$ref` — line-based (precise cursor targeting not yet implemented)
- [ ] **Rename** a symbol and update all its `$ref`s
- [ ] **Find usages / references**

## Known limitations

### Go to definition

- **Line-based targeting.** `$ref` is detected by parsing the single line under
the cursor, not the syntax node at the exact cursor position.
- **No JSON Pointer array indices.** A pointer segment is only resolved against
object keys. Refs whose path traverses an array — e.g.
`#/components/schemas/User/allOf/0` or
`#/paths/~1users/get/parameters/1/schema` — are not resolved
(`resolveSegment` only handles `MappingNode`, not `SequenceNode`).
- **No JSON Pointer escaping.** Pointer tokens are compared with object keys
verbatim, so the `~1` (`/`) and `~0` (`~`) escape sequences are not decoded.
This breaks refs into `paths`, whose keys are URL paths containing `/` — e.g.
`#/paths/~1users~1{id}` does not resolve to the key `/users/{id}`.
- **No percent-decoding.** The `$ref` fragment is treated as-is; percent-encoded
fragments (e.g. `%7B` for `{`) are not decoded before pointer resolution.
