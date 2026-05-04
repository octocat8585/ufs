# ufs

Unified File System (UFS) is a library that allows go apps to use the `fs.FS` interface to access file systems across many storage layers.

## Features

* Read-only `ReadFS`: Each file system implements Go's `fs.FS` interface with some extensions such as `Close() error` to manage the lifecycle of the access to the file system.
* Read-write `FS`: Provides read-write support for file systems that allow it.

## File Systems

Storage | URI         | Implementation | Description
------- | ----------- | -------------- | -----------
Null    | `nullfs://` | nullfs.go      | Acts as `/dev/null`. Any files written to this location are forgotten. Reads always resolve to empty files. All directories are read with no contents.
Memory  | `memfs://`  | memfs.go       | Stores files with virtual memory of the application. Forgotten when the application stops.
Local   | `file://`   | localfs.go     | Local files on the hard drive. The URI specifies the mount point of the local file system. Access to files outside of the mounted directory are not allowed. This is backed by Go's `os.RootFS`.
