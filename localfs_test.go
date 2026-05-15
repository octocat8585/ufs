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
	"testing"
)

const (
	testLocalFSName = "testing/testassets"
)

func TestLocalFSString(t *testing.T) {
	fsys, err := makeLocalFS(testLocalFSName)
	if err != nil {
		t.Fatal(err)
	}
	if got := fsys.String(); !strings.HasSuffix(got, testLocalFSName) {
		t.Errorf("String() should end with %q, got: %q", testLocalFSName, got)
	}
}

func TestIsLocalFSUri(t *testing.T) {
	testCases := []struct {
		name string
		want bool
	}{
		{name: "file:", want: true},
		{name: "file://", want: true},
		{name: "filefs://", want: false},
		{name: ".", want: true},
		{name: "/root/user", want: true},
		{name: "/tmp", want: true},
		{name: "mem://", want: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isLocalFSUri(tc.name)
			if got != tc.want {
				t.Errorf("got: %t, want: %t", got, tc.want)
			}
		})
	}
}

func TestLocalFS(t *testing.T) {
	dir := mustTemp(t)
	testFileSystem(t, newLocalFS, dir)
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

func mustLocalFS(tb testing.TB) *localFS {
	dir := mustTemp(tb)
	srcFS, err := makeLocalFS(testLocalFSName)
	if err != nil {
		tb.Fatalf("newLocalFS(%q) returned error, %s", testLocalFSName, err)
	}
	fsys, err := makeLocalFS(dir)
	if err != nil {
		tb.Fatalf("newLocalFS(%q) returned error, %s", dir, err)
	}
	if err := Rsync(srcFS, fsys, "."); err != nil {
		tb.Fatalf("cannot Rsync from %q to %q", testLocalFSName, dir)
	}
	return fsys
}
