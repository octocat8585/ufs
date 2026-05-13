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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xyproto/randomstring"
)

func testFileSystem(t *testing.T, newFSFunc func(name string) (FS, error), name string) {
	t.Helper()
	fsys := mustFS(t, newFSFunc, name)

	wantFiles := []string{"a", filepath.Join("ab", "b", "c"), filepath.Join("ab", "d", "c"), "def", "abc", "abc.txt", filepath.Join("temp", "abc.txt")}

	mkdirForTest(t, fsys, "ab", "b")
	mkdirForTest(t, fsys, "temp")
	mkdirForTest(t, fsys, "ab", "d")

	for _, name := range wantFiles {
		t.Run(fmt.Sprintf("crud_%s", name), func(t *testing.T) {
			wantData := randomString(1000)
			if wf, err := fsys.Create(name); err != nil {
				t.Errorf("cannot create file %q, %s", name, err)
			} else {
				info, err := wf.Stat()
				if err != nil {
					t.Errorf("cannot Stat() %q, %s", name, err)
				}
				if info == nil {
					t.Fatalf("info is nil")
				}
				if info.IsDir() != false {
					t.Errorf("%q is a directory, want file", name)
				}
				if n, err := io.WriteString(wf, wantData); err != nil {
					t.Errorf("cannot write file content to %q, %s", name, err)
				} else if n != len(wantData) {
					t.Errorf("contents written to file does not match the size got %d, want %d", n, len(wantData))
				}
				if err := wf.Close(); err != nil {
					t.Errorf("failed to Close() write file %q, %s", name, err)
				}
			}

			if rf, err := fsys.Open(name); err != nil {
				t.Errorf("cannot open file %q, %s", name, err)
			} else {
				if rf == nil {
					t.Fatal("rf is nil")
				}
				info, err := rf.Stat()
				if err != nil {
					t.Errorf("cannot Stat() %q, %s", name, err)
				}
				if info == nil {
					t.Fatal("info is nil")
				}
				if info.IsDir() != false {
					t.Errorf("%q is a directory, want file", name)
				}
				if got, err := io.ReadAll(rf); err != nil {
					t.Errorf("cannot read file content to %q, %s", name, err)
				} else if diff := cmp.Diff(wantData, string(got)); diff != "" {
					t.Errorf("io.ReadAll(%s) mismatch (-want +got):\n%s\nwant: %q\ngot: %q", name, diff, wantData, string(got))
				}
				if err := rf.Close(); err != nil {
					t.Errorf("failed to Close() read file %q, %s", name, err)
				}
			}
		})
	}

	if err := fstest.TestFS(fsys, wantFiles...); err != nil {
		t.Errorf("fstest.TestFS failed for %q: %v", name, err)
	}

	if err := fsys.Close(); err != nil {
		t.Errorf("error on Close(), %v", err)
	}
}

func mkdirForTest(tb testing.TB, fsys FS, dirs ...string) {
	tb.Helper()
	path := filepath.Join(dirs...)
	if err := fsys.MkdirAll(path, fs.ModePerm); err != nil {
		tb.Fatalf("cannot create directory %q, %s", path, err)
	}
}

func mustFS(tb testing.TB, newFSFunc func(name string) (FS, error), name string) FS {
	tb.Helper()

	fsys, err := newFSFunc(name)
	if err != nil {
		tb.Fatalf("FileSystem %q has an error, %s", name, err)
	}
	if fsys == nil {
		tb.Fatalf("FileSystem %q is nil", name)
	}

	return fsys
}

func randomString(size int) string {
	return randomstring.HumanFriendlyString(size)
}

func mustTemp(tb testing.TB) string {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			tb.Error(err)
		}
	})
	return tempDir
}

func mustTime(s string) time.Time {
	val, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return val
}

var (
	invalidPaths = []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}

	fsTestCases = []struct {
		name     string
		createFS func(tb testing.TB) (FS, error)
	}{
		{
			name: "localFS",
			createFS: func(tb testing.TB) (FS, error) {
				dir := mustTemp(tb)
				return newLocalFS(dir)
			},
		},
		{
			name: "tempMountFS",
			createFS: func(tb testing.TB) (FS, error) {
				return newTempMountFS("test://", func(string) error { return nil })
			},
		},
		{
			name: "memFS",
			createFS: func(tb testing.TB) (FS, error) {
				return newMemFS("mem://")
			},
		},
		{
			name: "nestFS",
			createFS: func(tb testing.TB) (FS, error) {
				return newNestFS("memory://")
			},
		},
	}

	readOnlyFSTestCases = []struct {
		name     string
		createFS func(tb testing.TB) (FS, error)
	}{
		{
			name: "nullFS",
			createFS: func(tb testing.TB) (FS, error) {
				return newNullFS("null://")
			},
		},
	}
)

func TestFSClose(t *testing.T) {
	t.Parallel()
	for _, tc := range append(fsTestCases, readOnlyFSTestCases...) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fsys, err := tc.createFS(t)
			if err != nil {
				t.Fatalf("cannot create file system: %s", err)
			}
			if err := fsys.Close(); err != nil {
				t.Errorf("Close() = %v, want nil", err)
			}
		})
	}
}

func TestFSMkdirAll(t *testing.T) {
	t.Parallel()
	for _, tc := range append(fsTestCases, readOnlyFSTestCases...) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fsys, err := tc.createFS(t)
			if err != nil {
				t.Fatalf("cannot create file system: %s", err)
			}
			defer fsys.Close()
			if err := fsys.MkdirAll("subdir", fs.ModePerm); err != nil {
				t.Errorf("MkdirAll() = %v, want nil", err)
			}
		})
	}
}

func TestFSReadFile(t *testing.T) {
	t.Parallel()
	for _, tc := range fsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			wantData := randomString(100)
			t.Parallel()
			fsys, err := tc.createFS(t)
			if err != nil {
				t.Fatalf("cannot create file system: %s", err)
			}
			defer fsys.Close()
			f, err := fsys.Create("readfile_test.txt")
			if err != nil {
				t.Fatalf("Create failed: %v", err)
			}
			if _, err := io.WriteString(f, wantData); err != nil {
				t.Fatalf("WriteString failed: %v", err)
			}
			if err := f.Close(); err != nil {
				t.Fatalf("Close failed: %v", err)
			}

			rfs, ok := fsys.(fs.ReadFileFS)
			if !ok {
				t.Skip("does not implement fs.ReadFileFS")
			}
			got, err := rfs.ReadFile("readfile_test.txt")
			if err != nil {
				t.Fatalf("ReadFile failed: %v", err)
			}
			if diff := cmp.Diff(wantData, string(got)); diff != "" {
				t.Errorf("ReadFile mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func verifyReadOnlyFS(t *testing.T, fsys fs.FS) {
	t.Helper()

	for _, path := range invalidPaths {
		t.Run("OpenInvalid", func(t *testing.T) {
			if _, err := fsys.Open(path); err == nil {
				t.Errorf("Open(%q) succeeded, want error", path)
			}
		})
	}
}

func verifyFS(t *testing.T, fsys FS) {
	verifyReadOnlyFS(t, fsys)

	for _, path := range invalidPaths {
		t.Run("CreateInvalid", func(t *testing.T) {
			if _, err := fsys.Create(path); err == nil {
				t.Errorf("Create(%q) succeeded, want error", path)
			}
		})
	}
}

func TestReadOnlyFS(t *testing.T) {
	t.Parallel()

	for _, tc := range readOnlyFSTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fsys, err := tc.createFS(t)
			if err != nil {
				t.Fatalf("cannot create file system, %s", err)
			}
			if fsys == nil {
				t.Fatalf("file sytem is nil")
			}
			defer func() {
				if err := fsys.Close(); err != nil {
					t.Errorf("second file system close() failed, %s", err)
				}
			}()
			verifyReadOnlyFS(t, fsys)
			if err := fsys.Close(); err != nil {
				t.Errorf("file system failed to close without errors, %s", err)
			}
		})
	}
}

func TestFS(t *testing.T) {
	t.Parallel()

	for _, tc := range fsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fsys, err := tc.createFS(t)
			if err != nil {
				t.Fatalf("cannot create file system, %s", err)
			}
			if fsys == nil {
				t.Fatalf("file sytem is nil")
			}
			defer func() {
				if err := fsys.Close(); err != nil {
					t.Errorf("second file system close() failed, %s", err)
				}
			}()
			verifyFS(t, fsys)
			if err := fsys.Close(); err != nil {
				t.Errorf("file system failed to close without errors, %s", err)
			}
		})
	}
}

func must(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Error(err)
	}
}
