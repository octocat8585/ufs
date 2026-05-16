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
	"strings"
)

const (
	angryFSPrefix = "angry:"
)

var (
	_ FS            = (*angryFS)(nil)
	_ fs.ReadFileFS = (*angryFS)(nil)
	_ fs.ReadDirFS  = (*angryFS)(nil)
	_ fs.ReadLinkFS = (*angryFS)(nil)
	_ fs.GlobFS     = (*angryFS)(nil)

	errAngry = fs.ErrInvalid
)

type angryFS struct {
	name string
}

func (fsys *angryFS) String() string {
	return fsys.name
}

func (fsys *angryFS) Open(name string) (fs.File, error) {
	if err := validPath("open", name); err != nil {
		return nil, err
	}
	return nil, errAngry
}

func (fsys *angryFS) Close() error {
	return errAngry
}

func (fsys *angryFS) Create(name string) (File, error) {
	if err := validPath("create", name); err != nil {
		return nil, err
	}

	return nil, errAngry
}

func (fsys *angryFS) MkdirAll(name string, perm fs.FileMode) error {
	if err := validPath("mkdir", name); err != nil {
		return err
	}
	return errAngry
}

func (fsys *angryFS) ReadFile(name string) ([]byte, error) {
	if err := validPath("readfile", name); err != nil {
		return nil, err
	}
	return nil, errAngry
}

func (fsys *angryFS) ReadLink(name string) (string, error) {
	if err := validPath("readlink", name); err != nil {
		return "", err
	}
	return "", errAngry
}

func (fsys *angryFS) Lstat(name string) (fs.FileInfo, error) {
	if err := validPath("lstat", name); err != nil {
		return nil, err
	}
	return nil, errAngry
}

func (fsys *angryFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := validPath("readdir", name); err != nil {
		return nil, err
	}
	return nil, errAngry
}

func (fsys *angryFS) Glob(pattern string) ([]string, error) {
	return nil, errAngry
}

func newAngryFS(name string) (*angryFS, error) {
	return makeAngryFS(name), nil
}

func makeAngryFS(name string) *angryFS {
	return &angryFS{
		name: name,
	}
}

func isAngryFSUri(name string) bool {
	return strings.HasPrefix(name, angryFSPrefix)
}
