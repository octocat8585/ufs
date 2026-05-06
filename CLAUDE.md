# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Test
go test ./...

# Test with race detector (required for presubmit)
CGO_ENABLED=1 go test -race ./...

# Run a single test
go test -run TestName ./...

# Lint (requires golangci-lint)
golangci-lint run

# Vet
go vet ./...

# Presubmit (lint + both test modes)
make presubmit
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

| File | Type | Status | Behavior |
|------|------|--------|----------|
| `nullfs.go` | `nullFS` | Implemented | `/dev/null` semantics — writes are discarded, reads return empty |
| `memfs.go` | `memFS` | Stub (TODO) | In-memory storage, lost when the process exits |
| `localfs.go` | `localFS` | Implemented | Local disk, mounted at a root path via `os.OpenRoot` (Go 1.24+); access outside mount point is disallowed |
| `angryfs.go` | `angryFS` | Implemented | Always returns `fs.ErrInvalid`; used to test error-handling paths |

### Supporting files

| File | Purpose |
|------|---------|
| `info.go` | `fsInfo` — concrete `fs.FileInfo` implementation |
| `util.go` | `validPath()` — validates paths against `fs.ValidPath` before all `Open`/`Create` calls |
| `testing_test.go` | `testFileSystem()` — shared CRUD + `fstest.TestFS` harness used by each backend's tests |

### Conventions

- Keep structs private; expose construction via `newXxxFS(name string) (FS, error)`
- Factory `name` arg follows a URI scheme: `null://`, `file:///abs/path` (localfs strips the `file://` prefix)
- All path operations call `validPath()` first — returns `*fs.PathError` for invalid paths
