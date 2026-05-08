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

package ufs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ File             = (*os.File)(nil)
	_ localFSInterface = (*localFS)(nil)
)

type localFSInterface interface {
	FS
	fs.ReadFileFS
	fs.ReadLinkFS
	// fs.ReadDirFS
	// fs.GlobFS
}

type localFS struct {
	osFS *os.Root
}

func (fsys *localFS) Open(name string) (fs.File, error) {
	if err := validPath("open", name); err != nil {
		return nil, err
	}
	return fsys.osFS.Open(name)
}

func (fsys *localFS) Close() error {
	return fsys.osFS.Close()
}

func (fsys *localFS) Create(name string) (File, error) {
	if err := validPath("create", name); err != nil {
		return nil, err
	}

	return fsys.osFS.Create(name)
}

func (fsys *localFS) MkdirAll(name string, perm fs.FileMode) error {
	return fsys.osFS.MkdirAll(name, perm)
}

func (fsys *localFS) ReadFile(name string) ([]byte, error) {
	if err := validPath("readfile", name); err != nil {
		return nil, err
	}
	return fsys.osFS.ReadFile(name)
}

func (fsys *localFS) ReadLink(name string) (string, error) {
	return fsys.osFS.Readlink(name)
}

func (fsys *localFS) Lstat(name string) (fs.FileInfo, error) {
	return fsys.osFS.Lstat(name)
}

func newLocalFS(name string) (FS, error) {
	name = strings.TrimPrefix(name, "file://")
	absPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	osFS, err := os.OpenRoot(absPath)
	if err != nil {
		return nil, err
	}
	return &localFS{
		osFS: osFS,
	}, nil
}
