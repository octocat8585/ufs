// Copyright 2026 Jeremy Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// ReadOnlyFS represents a read-only file system.
type ReadOnlyFS interface {
	fs.FS
	io.Closer
	// String returns a string description of the underly file system.
	String() string
}

type FS interface {
	ReadOnlyFS

	// Create a new file or directory within the file system.
	Create(name string) (File, error)

	// MkdirAll creates a directory.
	// If subdirectories do not exist within the chain they will also be created.
	MkdirAll(name string, perm fs.FileMode) error
}
