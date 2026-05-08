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

func TestNewNestFS(t *testing.T) {
	fsys, err := newNestFS("memfs://")
	if err != nil {
		t.Fatal(err)
	}
	if fsys == nil {
		t.Fatal("fsys is nil")
	}
}

func TestNestFS(t *testing.T) {
	testFileSystem(t, newNestFS, "memfs://")
}

func TestNestFSOpenInvalid(t *testing.T) {
	fsys, err := newNestFS("memfs://")
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}
	for _, path := range invalidPaths {
		if _, err := fsys.Open(path); err == nil {
			t.Errorf("Open(%q) succeeded, want error", path)
		}
	}
}

func TestNestFSCreateInvalid(t *testing.T) {
	fsys, err := newNestFS("memfs://")
	if err != nil {
		t.Fatal(err)
	}
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

func TestNestFSClose(t *testing.T) {
	fsys, err := newNestFS("memfs://")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}

func TestNestFSMkdirAll(t *testing.T) {
	fsys, err := newNestFS("memfs://")
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()
	if err := fsys.MkdirAll("subdir", fs.ModePerm); err != nil {
		t.Errorf("MkdirAll() = %v, want nil", err)
	}
}
