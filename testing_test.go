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

func verifyReadOnlyFS(tb testing.TB, createFS func(testing.TB) ReadFS) {
}

func verifyFS(tb testing.TB, createFS func(testing.TB) FS) {
	verifyReadOnlyFS(tb, func(iTB testing.TB) ReadFS {
		return createFS(iTB)
	})
}
