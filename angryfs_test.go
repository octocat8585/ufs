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
	"testing"
)

func TestIsAngryFSUri(t *testing.T) {
	testCases := []struct {
		name string
		want bool
	}{
		{
			name: "angry:",
			want: true,
		},
		{
			name: "angry://",
			want: true,
		},
		{
			name: "angryfs://",
			want: false,
		},
		{
			name: "mem://",
			want: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isAngryFSUri(tc.name)
			if got != tc.want {
				t.Errorf("got: %t, want: %t", got, tc.want)
			}
		})
	}
}

func TestNewAngryFS(t *testing.T) {
	fsys, err := newAngryFS("angry://")
	if err != nil {
		t.Fatal("newAngryFS() =", err)
	}
	if fsys == nil {
		t.Fatal("fsys is nil")
	}
}

func TestAngryFSOpen(t *testing.T) {
	fsys, err := newAngryFS("angry://")
	if err != nil {
		t.Fatal("newAngryFS() =", err)
	}

	// Test valid path - expects ErrInvalid as per current implementation
	_, err = fsys.Open("valid.txt")
	if err != fs.ErrInvalid {
		t.Errorf("Open(\"valid.txt\") = %v, want %v", err, fs.ErrInvalid)
	}

	// Test invalid paths
	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}

	for _, path := range invalidPaths {
		_, err := fsys.Open(path)
		if err == nil {
			t.Errorf("Open(%q) succeeded, want error", path)
		} else if _, ok := err.(*fs.PathError); !ok {
			t.Errorf("Open(%q) returned %T, want *fs.PathError", path, err)
		}
	}
}

func TestAngryFSClose(t *testing.T) {
	fsys, err := newAngryFS("angry://")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.Close(); err != fs.ErrInvalid {
		t.Errorf("Close() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSMkdirAll(t *testing.T) {
	fsys, err := newAngryFS("angry://")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.MkdirAll("dir", fs.ModePerm); err != fs.ErrInvalid {
		t.Errorf("MkdirAll() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSReadFile(t *testing.T) {
	afs := makeAngryFS()
	if _, err := afs.ReadFile("foo.txt"); err != fs.ErrInvalid {
		t.Errorf("ReadFile() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSReadLink(t *testing.T) {
	afs := makeAngryFS()
	if _, err := afs.ReadLink("foo.txt"); err != fs.ErrInvalid {
		t.Errorf("ReadLink() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSLstat(t *testing.T) {
	afs := makeAngryFS()
	if _, err := afs.Lstat("foo.txt"); err != fs.ErrInvalid {
		t.Errorf("Lstat() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSReadDir(t *testing.T) {
	afs := makeAngryFS()
	if _, err := afs.ReadDir("dir"); err != fs.ErrInvalid {
		t.Errorf("ReadDir() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSGlob(t *testing.T) {
	afs := makeAngryFS()
	if _, err := afs.Glob("*.txt"); err != fs.ErrInvalid {
		t.Errorf("Glob() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSCreate(t *testing.T) {
	fsys, err := newAngryFS("angry://")
	if err != nil {
		t.Fatal("newAngryFS() =", err)
	}

	// Test valid path - expects ErrInvalid
	_, err = fsys.Create("valid.txt")
	if err != fs.ErrInvalid {
		t.Errorf("Create(\"valid.txt\") = %v, want %v", err, fs.ErrInvalid)
	}

	// Test invalid paths - checking for PathError even though implementation might fail early
	// But `angryFS` implementation checks ValidPath first, so it should be PathError.
	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}

	for _, path := range invalidPaths {
		_, err := fsys.Create(path)
		if err == nil {
			t.Errorf("Create(%q) succeeded, want error", path)
		}
	}
}
