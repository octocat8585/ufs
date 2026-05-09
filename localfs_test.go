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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLocalFS(t *testing.T) {
	dir := mustTemp(t)
	testFileSystem(t, newLocalFS, dir)
}

func TestLocalFSCreateInvalid(t *testing.T) {
	dir := mustTemp(t)
	fsys := mustFS(t, newLocalFS, dir)
	defer fsys.Close()

	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}
	for _, path := range invalidPaths {
		if _, err := fsys.Create(path); err == nil {
			t.Errorf("Create(%q) succeeded, want error", path)
		}
	}
}

func TestLocalFSReadFile(t *testing.T) {
	dir := mustTemp(t)
	fsys := mustFS(t, newLocalFS, dir)
	defer fsys.Close()

	wantData := randomString(100)

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
		t.Fatal("localFS does not implement fs.ReadFileFS")
	}
	got, err := rfs.ReadFile("readfile_test.txt")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if diff := cmp.Diff(wantData, string(got)); diff != "" {
		t.Errorf("ReadFile mismatch (-want +got):\n%s", diff)
	}
}

func TestLocalFSLstat(t *testing.T) {
	dir := mustTemp(t)
	fsys := mustFS(t, newLocalFS, dir)
	defer fsys.Close()

	f, err := fsys.Create("lstat_file.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()

	if err := os.Symlink("lstat_file.txt", filepath.Join(dir, "lstat_link.txt")); err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	lfs, ok := fsys.(fs.ReadLinkFS)
	if !ok {
		t.Fatal("localFS does not implement fs.ReadLinkFS")
	}

	t.Run("regular_file", func(t *testing.T) {
		info, err := lfs.Lstat("lstat_file.txt")
		if err != nil {
			t.Fatalf("Lstat failed: %v", err)
		}
		if info.IsDir() {
			t.Error("IsDir() = true, want false")
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			t.Errorf("Mode() has ModeSymlink set on a regular file: %v", info.Mode())
		}
	})

	t.Run("symlink", func(t *testing.T) {
		info, err := lfs.Lstat("lstat_link.txt")
		if err != nil {
			t.Fatalf("Lstat on symlink failed: %v", err)
		}
		if info.Mode()&fs.ModeSymlink == 0 {
			t.Errorf("Mode() missing ModeSymlink for symlink: %v", info.Mode())
		}
	})
}

func TestLocalFSReadLink(t *testing.T) {
	dir := mustTemp(t)
	fsys := mustFS(t, newLocalFS, dir)
	defer fsys.Close()

	f, err := fsys.Create("target.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()

	if err := os.Symlink("target.txt", filepath.Join(dir, "link.txt")); err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	lfs, ok := fsys.(fs.ReadLinkFS)
	if !ok {
		t.Fatal("localFS does not implement fs.ReadLinkFS")
	}

	got, err := lfs.ReadLink("link.txt")
	if err != nil {
		t.Fatalf("ReadLink failed: %v", err)
	}
	if got != "target.txt" {
		t.Errorf("ReadLink = %q, want %q", got, "target.txt")
	}
}
