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
	"testing"
)

func TestNewNullFS(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}
	if fsys == nil {
		t.Fatal("fsys is nil")
	}

	f, err := fsys.Open("foo.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	nf, ok := f.(io.Writer)
	if !ok {
		t.Fatal("file does not implement Write")
	}

	data := []byte("hello world")
	n, err := nf.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(data) {
		t.Errorf("Write() = %d, want %d", n, len(data))
	}
}

func TestNullFSOpenInvalid(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}
	invalidPaths := []string{
		"/absolute/path",
		"../relative/path",
		"invalid/../path",
	}

	for _, path := range invalidPaths {
		_, err := fsys.Open(path)
		if err == nil {
			t.Errorf("Open(%q) succeeded, want error", path)
		}
	}
}

func TestNullFSCreate(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}

	// Valid Create
	f, err := fsys.Create("created.txt")
	if err != nil {
		t.Fatalf("Create(\"created.txt\") failed: %v", err)
	}
	defer f.Close()

	if f == nil {
		t.Fatal("Created file is nil")
	}

	// Invalid Create
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

func TestNullFileOperations(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Open("testops.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Read
	b := make([]byte, 10)
	n, err := f.Read(b)
	if n != 0 {
		t.Errorf("Read() = %d, want 0", n)
	}
	if err != io.EOF {
		t.Errorf("Read() error = %v, want io.EOF", err)
	}

	// Close
	if err := f.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}

	// Stat
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat() failed: %v", err)
	}
	if info.Name() != "testops.txt" {
		t.Errorf("Name() = %q, want %q", info.Name(), "testops.txt")
	}
	if info.Size() != 0 {
		t.Errorf("Size() = %d, want 0", info.Size())
	}
	if info.IsDir() {
		t.Error("IsDir() = true, want false")
	}
	if info.Mode() != fs.ModePerm {
		t.Errorf("Mode() = %v, want %v", info.Mode(), fs.ModePerm)
	}
}
