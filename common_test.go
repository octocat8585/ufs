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
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
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

func TestFSConventions(t *testing.T) {
	srcFS, err := newLocalFS(testLocalFSName)
	if err != nil {
		t.Fatalf("cannot mount localFS(%q), %s", testLocalFSName, err)
	}
	for _, fsysTC := range getReadWriteTestCaseList() {
		t.Run(fsysTC.name, func(t *testing.T) {
			t.Parallel()
			fsys := fsysTC.createFS(t)
			if err := Rsync(srcFS, fsys, cwdPath); err != nil {
				t.Errorf("rsync failed with error, %s", err)
			}

			allFilenames, err := List(srcFS, cwdPath)
			if err != nil {
				t.Fatal(err)
			}
			if len(allFilenames) == 0 {
				t.Fatal("expected at least 1 file name")
			}
			if err := fstest.TestFS(fsys, allFilenames...); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestFSClose(t *testing.T) {
	// TODO: Change to getAllTestCaseList
	for _, fsysTC := range getReadWriteTestCaseList() {
		t.Run(fsysTC.name, func(t *testing.T) {
			t.Parallel()
			fsys := fsysTC.createFS(t)
			for i := range 10 {
				if err := fsys.Close(); err != nil {
					t.Errorf("Close() [%d] failed with error, %s", i, err)
				}
			}

			t.Run("Open", func(t *testing.T) {
				f, err := fsys.Open(".")
				if f != nil {
					buf := make([]byte, 64)
					if bytesWritten, err := f.Read(buf); bytesWritten != 0 && err != fs.ErrClosed {
						t.Errorf("Open('.') worked for a closed file system, read %d (want: 0) bytes with error= %s (want: fs.ErrClosed), file: %+v", bytesWritten, err, f)
					}
				}
				if !errors.Is(err, fs.ErrClosed) {
					t.Errorf("Open('.') did not return fs.ErrClosed, got: %s", err)
				}
			})

			t.Run("Create", func(t *testing.T) {
				f, err := fsys.Create("file.txt")
				if f != nil {
					buf := make([]byte, 64)
					if bytesWritten, err := f.Write(buf); bytesWritten != 0 && err != fs.ErrClosed {
						t.Errorf("Create('file.txt') worked for a closed file system, read %d (want: 0) bytes with error= %s (want: fs.ErrClosed), file: %+v", bytesWritten, err, f)
					}
				}
				if !errors.Is(err, fs.ErrClosed) {
					t.Errorf("Create('file.txt') did not return fs.ErrClosed, got: %s", err)
				}
			})

			t.Run("MkdirAll", func(t *testing.T) {
				if err := fsys.MkdirAll("a/b/c", fs.ModePerm); !errors.Is(err, fs.ErrClosed) {
					t.Errorf("MkdirAll('a/b/c') did not return fs.ErrClosed, got: %s", err)
				}
			})

			t.Run("ReadDir", func(t *testing.T) {
				entries, err := fs.ReadDir(fsys, ".")
				if len(entries) != 0 {
					t.Errorf("ReadDir('.') returned results for a closed file system, got: %v", entries)
				}
				if !errors.Is(err, fs.ErrClosed) {
					t.Errorf("ReadDir('.') did not return fs.ErrClosed, got: %s", err)
				}
			})
		})
	}
}

func TestFSString(t *testing.T) {
	for _, tc := range getAllTestCaseList() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fsys := tc.createFS(t)
			if got := fsys.String(); !strings.Contains(got, tc.wantString) {
				t.Errorf("%s.String() should contain %q: got: %q", fsys, tc.wantString, got)
			}
		})
	}
}
