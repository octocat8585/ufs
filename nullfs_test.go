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

func TestNullFSClose(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}

func TestNullFSMkdirAll(t *testing.T) {
	fsys, err := newNullFS("null://")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.MkdirAll("dir", fs.ModePerm); err != nil {
		t.Errorf("MkdirAll() = %v, want nil", err)
	}
}

func TestNullFSReadFile(t *testing.T) {
	nfs := &nullFS{}
	got, err := nfs.ReadFile("foo.txt")
	if err != nil {
		t.Fatalf("ReadFile() = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Errorf("ReadFile() = %q, want empty", got)
	}
}

func TestNullFSReadLink(t *testing.T) {
	nfs := &nullFS{}
	got, err := nfs.ReadLink("link.txt")
	if err != nil {
		t.Fatalf("ReadLink() = %v, want nil", err)
	}
	if got != "" {
		t.Errorf("ReadLink() = %q, want empty string", got)
	}
}

func TestNullFSLstat(t *testing.T) {
	nfs := &nullFS{}

	t.Run("file", func(t *testing.T) {
		info, err := nfs.Lstat("foo.txt")
		if err != nil {
			t.Fatalf("Lstat() = %v, want nil", err)
		}
		if info.IsDir() {
			t.Error("IsDir() = true, want false")
		}
	})

	t.Run("dir", func(t *testing.T) {
		info, err := nfs.Lstat("dir/")
		if err != nil {
			t.Fatalf("Lstat() = %v, want nil", err)
		}
		if !info.IsDir() {
			t.Error("IsDir() = false, want true")
		}
		if info.Mode()&fs.ModeDir == 0 {
			t.Errorf("Mode() missing ModeDir: %v", info.Mode())
		}
	})
}

func TestNullFSReadDir(t *testing.T) {
	nfs := &nullFS{}
	entries, err := nfs.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir() = %v, want nil", err)
	}
	if len(entries) != 0 {
		t.Errorf("ReadDir() = %d entries, want 0", len(entries))
	}
}

func TestNullFSGlob(t *testing.T) {
	nfs := &nullFS{}
	matches, err := nfs.Glob("*.txt")
	if err != nil {
		t.Fatalf("Glob() = %v, want nil", err)
	}
	if len(matches) != 0 {
		t.Errorf("Glob() = %d matches, want 0", len(matches))
	}
}

func TestNullFileStatDir(t *testing.T) {
	f := newNullFile("subdir/")
	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Error("IsDir() = false, want true for trailing-slash name")
	}
	if info.Mode()&fs.ModeDir == 0 {
		t.Errorf("Mode() missing ModeDir: %v", info.Mode())
	}
}

func TestNullFileSeek(t *testing.T) {
	f := newNullFile("seek.txt")
	pos, err := f.Seek(100, io.SeekStart)
	if err != nil {
		t.Fatalf("Seek() = %v, want nil", err)
	}
	if pos != 0 {
		t.Errorf("Seek() = %d, want 0", pos)
	}
}

func TestNullFileReadAt(t *testing.T) {
	f := newNullFile("readat.txt")
	buf := make([]byte, 10)
	n, err := f.ReadAt(buf, 0)
	if err != io.EOF {
		t.Errorf("ReadAt() error = %v, want io.EOF", err)
	}
	if n != 0 {
		t.Errorf("ReadAt() = %d, want 0", n)
	}
}

func TestNullFileWriteString(t *testing.T) {
	f := newNullFile("write.txt")
	n, err := f.WriteString("hello world")
	if err != nil {
		t.Fatalf("WriteString() = %v, want nil", err)
	}
	if n != 11 {
		t.Errorf("WriteString() = %d, want 11", n)
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
