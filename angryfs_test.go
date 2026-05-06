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
