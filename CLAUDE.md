# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Test
go test ./...

# Run a single test
go test -run TestName ./...

# Lint (requires golangci-lint)
golangci-lint run

# Vet
go vet ./...
```

## Architecture

`ufs` is a Go library providing a unified virtual file system abstraction. The module is `github.com/jeremyje/ufs`.

### Core interfaces (`ufs.go`)

- `ReadFile` — read-only file, wraps `fs.File`
- `File` — read-write file, extends `ReadFile` with `io.ReaderAt`, `io.ReadWriteSeeker`, and `io.StringWriter`
- `ReadFS` — read-only file system, wraps `fs.FS` and adds `io.Closer` for lifecycle management
- `FS` — read-write file system, extends `ReadFS` with `Create(name string) (File, error)` and `MkdirAll(name string, perm fs.FileMode) error`

### Implementations

Each file system backend is a private struct in its own file implementing the `FS` interface:

| File | Type | Behavior |
|------|------|----------|
| `nullfs.go` | `nullFS` | `/dev/null` semantics — writes are discarded, reads return empty |
| `memfs.go` | `memFS` | In-memory storage, lost when the process exits |
| `localfs.go` | `localFS` | Local disk, mounted at a root path via `os.RootFS`; access outside the mount point is disallowed |

All three are currently stubs with `TODO` comments. When implementing, keep the struct private and expose construction through a package-level factory function (e.g., `NewNullFS() FS`).
