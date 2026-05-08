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
)

var (
	_ FS = (*nestFS)(nil)
)

// nestFS is a wrapper for a base FS that supports automatic mounting of archives.
// This means that any archive can opened and read automatically. Archives are revealed via $filename.d name pattern.
type nestFS struct {
	fsys FS
}

func (fsys *nestFS) Open(name string) (fs.File, error) {
	if err := validPath("open", name); err != nil {
		return nil, err
	}

	return fsys.fsys.Open(name)
}

func (fsys *nestFS) Close() error {
	return fsys.fsys.Close()
}

func (fsys *nestFS) Create(name string) (File, error) {
	if err := validPath("create", name); err != nil {
		return nil, err
	}

	// TODO: 
	return fsys.fsys.Create(name)
}

func (fsys *nestFS) MkdirAll(name string, perm fs.FileMode) error {
	// TODO: 
	return fsys.fsys.MkdirAll(name, perm)
}

func newNestFS(name string) (FS, error) {
	fsys, err := newBaseFS(name)
	if err != nil {
		return nil, err
	}
	return &nestFS{
		fsys: fsys,
	}, nil
}
