// Package ufs is a virtual system library for Go. It allows go applications to mount different storage systems represented as file systems to be used in applications.
package ufs

import (
	"io"
	"io/fs"
)

// ReadFile represents a read-only file.
type ReadFile interface {
	fs.File
}

// File represents a read-write file.
type File interface {
	ReadFile
	io.ReaderAt
	io.ReadWriteSeeker
	io.StringWriter
}

// ReadFS represents a read-only file system.
type ReadFS interface {
	fs.FS
	io.Closer
}

type FS interface {
	ReadFS

	// Create a new file or directory within the file system.
	Create(name string) (File, error)

	// MkdirAll creates a directory.
	// If subdirectories do not exist within the chain they will also be created.
	MkdirAll(name string, perm fs.FileMode) error
}
