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
	"fmt"
	"io/fs"
	"reflect"
	"testing"
)

func TestInvalidPath(t *testing.T) {
	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
		"",
	}

	for _, fsysTC := range getAllTestCaseList() {
		for _, path := range invalidPaths {
			t.Run(fmt.Sprintf("Open/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				_, err := fsys.Open(path)
				assertInvalidPathError(t, path, err, "open")
			})

			t.Run(fmt.Sprintf("ReadDir/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				_, err := fsys.ReadDir(path)
				assertInvalidPathError(t, path, err, "readdir")
			})

			t.Run(fmt.Sprintf("Create/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				_, err := fsys.Create(path)
				assertInvalidPathError(t, path, err, "create")
			})

			t.Run(fmt.Sprintf("MkdirAll/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				err := fsys.MkdirAll(path, fs.ModeDir)
				assertInvalidPathError(t, path, err, "mkdir")
			})

			t.Run(fmt.Sprintf("ReadFileFS/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				if rf, ok := fsys.(fs.ReadFileFS); ok {
					_, err := rf.ReadFile(path)
					assertInvalidPathError(t, path, err, "readfile")
				}
			})

			t.Run(fmt.Sprintf("ReadLink/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				if rf, ok := fsys.(fs.ReadLinkFS); ok {
					_, err := rf.ReadLink(path)
					assertInvalidPathError(t, path, err, "readlink")
				}
			})

			t.Run(fmt.Sprintf("Lstat/%s/%s", fsysTC.name, path), func(t *testing.T) {
				t.Parallel()
				fsys := fsysTC.createFS(t)
				if rf, ok := fsys.(fs.ReadLinkFS); ok {
					_, err := rf.Lstat(path)
					assertInvalidPathError(t, path, err, "lstat")
				}
			})
		}
	}
}

func assertInvalidPathError(t *testing.T, path string, err error, wantOp string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s(%q) succeeded, want error", wantOp, path)
		return
	}
	if perr, ok := err.(*fs.PathError); ok {
		if wantOp != perr.Op {
			t.Errorf("fs.PathError.Op mismatch, got: %q, want: %q", perr.Op, wantOp)
		}
		if path != perr.Path {
			t.Errorf("fs.PathError.Path mismatch, got: %q, want: %q", perr.Path, path)
		}
	} else {
		t.Errorf("%q is not a *fs.PathError, got: %q", err, reflect.TypeOf(err).Name())
	}
}
