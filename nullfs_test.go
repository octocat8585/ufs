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

func TestNullFSReadFile(t *testing.T) {
	nfs := makeNullFS()
	got, err := nfs.ReadFile("foo.txt")
	if err != nil {
		t.Fatalf("ReadFile() = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Errorf("ReadFile() = %q, want empty", got)
	}
}

func TestNullFSReadLink(t *testing.T) {
	nfs := makeNullFS()
	_, err := nfs.ReadLink("link.txt")
	if err == nil {
		t.Fatal("ReadLink() = nil error, want error (nullFS has no symlinks)")
	}
}

func TestNullFSLstat(t *testing.T) {
	nfs := makeNullFS()

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
		// "dir" is a valid fs.FS path; nullFS has no storage so Lstat returns file info.
		info, err := nfs.Lstat("dir")
		if err != nil {
			t.Fatalf("Lstat() = %v, want nil", err)
		}
		if info == nil {
			t.Fatal("Lstat() returned nil info")
		}
	})
}

func TestNullFSReadDir(t *testing.T) {
	nfs := makeNullFS()
	entries, err := nfs.ReadDir(cwdPath)
	if err != nil {
		t.Fatalf("ReadDir() = %v, want nil", err)
	}
	if len(entries) != 0 {
		t.Errorf("ReadDir() = %d entries, want 0", len(entries))
	}
	f, err := nfs.Open(cwdPath)
	if err != nil {
		t.Errorf("Open('.') returned error, %s", err)
	}
	rdf, ok := f.(fs.ReadDirFile)
	if ok {
		if readBytes, err := rdf.Read(nil); err == nil || readBytes != 0 {
			t.Errorf("fs.ReadDirFile.Read() want: (0, error) got: (%d, %s)", readBytes, err)
		}
		if dirs, err := rdf.ReadDir(-1); err != nil {
			t.Errorf("fs.ReadDirFile.ReadDir(-1) returned error, %s", err)
		} else if len(dirs) != 0 {
			t.Errorf("fs.ReadDirFile.ReadDir(-1) want: [], got: %v", dirs)
		}
		if stat, err := rdf.Stat(); err != nil {
			t.Errorf("fs.ReadDirFile.Stat() returned error, %s", err)
		} else if stat != nullDirStat {
			t.Errorf("fs.ReadDirFile.Stat() want: %v, got: %v", nullDirStat, stat)
		}
		if err := rdf.Close(); err != nil {
			t.Errorf("fs.ReadDirFile returned error, %s", err)
		}
	} else {
		t.Errorf("%+v is not a fs.ReadDirFile", rdf)
	}
}

func TestNullFSGlob(t *testing.T) {
	nfs := makeNullFS()
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
