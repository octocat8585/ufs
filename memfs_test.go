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

func TestNewMemFS(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	if fsys == nil {
		t.Fatal("fsys is nil")
	}
}

func TestMemFS(t *testing.T) {
	testFileSystem(t, newMemFS, "memory://test")
}

func TestMemFSOpenInvalid(t *testing.T) {
	fsys, err := newMemFS("memory://test")
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

func TestMemFSCreate(t *testing.T) {
	fsys, err := newMemFS("memory://test")
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

func TestMemFileOperations(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	// Create and write
	f, err := fsys.Create("testops.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Write
	n, err := f.Write([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 11 {
		t.Errorf("Write() = %d, want 11", n)
	}

	// Seek to beginning
	pos, err := f.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 0 {
		t.Errorf("Seek(0, Start) = %d, want 0", pos)
	}

	// Read
	buf := make([]byte, 11)
	n, err = f.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 11 {
		t.Errorf("Read() = %d, want 11", n)
	}
	if string(buf) != "hello world" {
		t.Errorf("Read() = %q, want %q", string(buf), "hello world")
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
	if info.Size() != 11 {
		t.Errorf("Size() = %d, want 11", info.Size())
	}
	if info.IsDir() {
		t.Error("IsDir() = true, want false")
	}
	if info.Mode() != fs.ModePerm {
		t.Errorf("Mode() = %v, want %v", info.Mode(), fs.ModePerm)
	}

	// Read after EOF
	n, err = f.Read(buf)
	if n != 0 {
		t.Errorf("Read() after EOF = %d, want 0", n)
	}
	if err != io.EOF {
		t.Errorf("Read() error after EOF = %v, want io.EOF", err)
	}
}

func TestMemFileSeek(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	f, err := fsys.Create("seek.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("abcdefghijklm")

	// Seek to middle
	pos, err := f.Seek(5, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 5 {
		t.Errorf("Seek(5, Start) = %d, want 5", pos)
	}

	buf := make([]byte, 3)
	n, err := f.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("Read() = %d, want 3", n)
	}
	if string(buf) != "fgh" {
		t.Errorf("Read() = %q, want %q", string(buf), "fgh")
	}

	// SeekCurrent
	pos, err = f.Seek(2, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 10 {
		t.Errorf("Seek(2, Current) = %d, want 10", pos)
	}

	// SeekEnd
	pos, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 13 {
		t.Errorf("Seek(0, End) = %d, want 13", pos)
	}

	// Invalid whence
	_, err = f.Seek(0, 99)
	if err == nil {
		t.Error("Seek(0, 99) succeeded, want error")
	}
}

func TestMemFileReadAt(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	f, err := fsys.Create("readat.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello world")

	// Read from middle (should not hit EOF)
	buf := make([]byte, 5)
	n, err := f.ReadAt(buf, 1)
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("ReadAt() = %d, want 5", n)
	}
	if string(buf) != "ello " {
		t.Errorf("ReadAt() = %q, want %q", string(buf), "ello ")
	}

	// Read to end (should return EOF)
	buf2 := make([]byte, 5)
	n, err = f.ReadAt(buf2, 6)
	if err != io.EOF {
		t.Errorf("ReadAt() at end err = %v, want io.EOF", err)
	}
	if n != 5 {
		t.Errorf("ReadAt() at end = %d, want 5", n)
	}
	if string(buf2) != "world" {
		t.Errorf("ReadAt() at end = %q, want %q", string(buf2), "world")
	}

	// Read beyond end
	buf3 := make([]byte, 5)
	n, err = f.ReadAt(buf3, 11)
	if err != io.EOF {
		t.Errorf("ReadAt() beyond end err = %v, want io.EOF", err)
	}
	if n != 0 {
		t.Errorf("ReadAt() beyond end = %d, want 0", n)
	}
}

func TestMemFSDirectory(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	// Create directory
	err = fsys.MkdirAll("subdir", fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// Open directory
	dir, err := fsys.Open("subdir")
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()

	// Stat directory
	info, err := dir.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Error("IsDir() = false, want true")
	}

	// Readdir on file should fail
	f, err := fsys.Create("subdir/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if mf, ok := f.(*memFile); ok {
		if _, err := mf.Readdir(1); err == nil {
			t.Error("Readdir on file succeeded, want error")
		}
	} else {
		t.Fatal("file is not *memFile")
	}
}

func TestMemFSFilePersistence(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	// Create and write a file
	f, err := fsys.Create("persist.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("persistent data")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Reopen and read
	f2, err := fsys.Open("persist.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	data, err := io.ReadAll(f2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "persistent data" {
		t.Errorf("ReadAll() = %q, want %q", string(data), "persistent data")
	}
}

func TestMemFSReadFile(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	f, err := fsys.Create("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello world")
	f.Close()

	rfs := fsys.(fs.ReadFileFS)

	t.Run("valid", func(t *testing.T) {
		got, err := rfs.ReadFile("hello.txt")
		if err != nil {
			t.Fatalf("ReadFile() = %v, want nil", err)
		}
		if string(got) != "hello world" {
			t.Errorf("ReadFile() = %q, want %q", got, "hello world")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		if _, err := rfs.ReadFile("missing.txt"); err == nil {
			t.Error("ReadFile(missing) succeeded, want error")
		}
	})

	t.Run("invalid_path", func(t *testing.T) {
		if _, err := rfs.ReadFile("../escape.txt"); err == nil {
			t.Error("ReadFile(../escape.txt) succeeded, want error")
		}
	})
}

func TestMemFSReadLink(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := fsys.Create("file.txt")
	f.Close()

	lfs := fsys.(fs.ReadLinkFS)

	t.Run("existing_file_not_a_symlink", func(t *testing.T) {
		if _, err := lfs.ReadLink("file.txt"); err == nil {
			t.Error("ReadLink on regular file succeeded, want error")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		if _, err := lfs.ReadLink("missing.txt"); err == nil {
			t.Error("ReadLink(missing) succeeded, want error")
		}
	})

	t.Run("invalid_path", func(t *testing.T) {
		if _, err := lfs.ReadLink("../escape.txt"); err == nil {
			t.Error("ReadLink(../escape) succeeded, want error")
		}
	})
}

func TestMemFSLstat(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	fsys.MkdirAll("mydir", fs.ModePerm)
	f, _ := fsys.Create("mydir/file.txt")
	f.WriteString("data")
	f.Close()

	lfs := fsys.(fs.ReadLinkFS)

	t.Run("regular_file", func(t *testing.T) {
		info, err := lfs.Lstat("mydir/file.txt")
		if err != nil {
			t.Fatalf("Lstat() = %v, want nil", err)
		}
		if info.Name() != "file.txt" {
			t.Errorf("Name() = %q, want %q", info.Name(), "file.txt")
		}
		if info.Size() != 4 {
			t.Errorf("Size() = %d, want 4", info.Size())
		}
		if info.IsDir() {
			t.Error("IsDir() = true, want false")
		}
	})

	t.Run("directory", func(t *testing.T) {
		info, err := lfs.Lstat("mydir")
		if err != nil {
			t.Fatalf("Lstat() = %v, want nil", err)
		}
		if !info.IsDir() {
			t.Error("IsDir() = false, want true for directory")
		}
		if info.Mode()&fs.ModeDir == 0 {
			t.Errorf("Mode() missing ModeDir: %v", info.Mode())
		}
	})

	t.Run("root", func(t *testing.T) {
		info, err := lfs.Lstat(".")
		if err != nil {
			t.Fatalf("Lstat(.) = %v, want nil", err)
		}
		if !info.IsDir() {
			t.Error("IsDir() = false, want true for root")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		if _, err := lfs.Lstat("missing.txt"); err == nil {
			t.Error("Lstat(missing) succeeded, want error")
		}
	})

	t.Run("invalid_path", func(t *testing.T) {
		if _, err := lfs.Lstat("../escape"); err == nil {
			t.Error("Lstat(../escape) succeeded, want error")
		}
	})
}

func TestMemFSReadDir(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	fsys.MkdirAll("docs", fs.ModePerm)
	for _, name := range []string{"a.txt", "b.txt"} {
		f, _ := fsys.Create("docs/" + name)
		f.Close()
	}

	dfs := fsys.(fs.ReadDirFS)

	t.Run("populated_dir", func(t *testing.T) {
		entries, err := dfs.ReadDir("docs")
		if err != nil {
			t.Fatalf("ReadDir() = %v, want nil", err)
		}
		if len(entries) != 2 {
			t.Errorf("ReadDir() = %d entries, want 2", len(entries))
		}
	})

	t.Run("root", func(t *testing.T) {
		entries, err := dfs.ReadDir(".")
		if err != nil {
			t.Fatalf("ReadDir(.) = %v, want nil", err)
		}
		if len(entries) == 0 {
			t.Error("ReadDir(.) returned 0 entries, want at least 1")
		}
	})

	t.Run("on_file", func(t *testing.T) {
		f, _ := fsys.Create("plain.txt")
		f.Close()
		if _, err := dfs.ReadDir("plain.txt"); err == nil {
			t.Error("ReadDir on a file succeeded, want error")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		if _, err := dfs.ReadDir("missing"); err == nil {
			t.Error("ReadDir(missing) succeeded, want error")
		}
	})

	t.Run("invalid_path", func(t *testing.T) {
		if _, err := dfs.ReadDir("../escape"); err == nil {
			t.Error("ReadDir(../escape) succeeded, want error")
		}
	})
}

func TestMemFSGlob(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	fsys.MkdirAll("src", fs.ModePerm)
	for _, name := range []string{"src/foo.go", "src/bar.go", "src/README.md"} {
		f, _ := fsys.Create(name)
		f.Close()
	}

	gfs := fsys.(fs.GlobFS)

	t.Run("match_go_files", func(t *testing.T) {
		matches, err := gfs.Glob("src/*.go")
		if err != nil {
			t.Fatalf("Glob() = %v, want nil", err)
		}
		if len(matches) != 2 {
			t.Errorf("Glob(src/*.go) = %v, want 2 matches", matches)
		}
	})

	t.Run("match_all_in_src", func(t *testing.T) {
		matches, err := gfs.Glob("src/*")
		if err != nil {
			t.Fatalf("Glob() = %v, want nil", err)
		}
		if len(matches) != 3 {
			t.Errorf("Glob(src/*) = %v, want 3 matches", matches)
		}
	})

	t.Run("no_match", func(t *testing.T) {
		matches, err := gfs.Glob("src/*.xyz")
		if err != nil {
			t.Fatalf("Glob() = %v, want nil", err)
		}
		if len(matches) != 0 {
			t.Errorf("Glob(src/*.xyz) = %v, want 0 matches", matches)
		}
	})

	t.Run("invalid_pattern", func(t *testing.T) {
		if _, err := gfs.Glob("[invalid"); err == nil {
			t.Error("Glob([invalid) succeeded, want error")
		}
	})
}

func TestMemFSReaddirAll(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.MkdirAll("parent", fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		f, err := fsys.Create("parent/" + name)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	dir, err := fsys.Open("parent")
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()

	entries, err := dir.(*memFile).Readdir(-1)
	if err != nil {
		t.Fatalf("Readdir(-1) = %v, want nil", err)
	}
	if len(entries) != 2 {
		t.Errorf("Readdir(-1) = %d entries, want 2", len(entries))
	}
}

func TestMemFSReaddirPaginated(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.MkdirAll("paged", fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		f, err := fsys.Create("paged/" + name)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	dir, err := fsys.Open("paged")
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	mf := dir.(*memFile)

	for i := range 3 {
		e, err := mf.Readdir(1)
		if err != nil || len(e) != 1 {
			t.Fatalf("Readdir(1) call %d: got %d entries, err=%v", i+1, len(e), err)
		}
	}
	// Exhausted: next call with n>0 must return io.EOF
	if _, err := mf.Readdir(1); err != io.EOF {
		t.Errorf("Readdir(1) after exhaustion = %v, want io.EOF", err)
	}
}

func TestMemFileReadDirOnFile(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Create("regular.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.(*memFile).ReadDir(0); err == nil {
		t.Error("ReadDir on a regular file succeeded, want error")
	}
}

func TestMemFileSeekNegative(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fsys.Create("neg.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Seek(-1, io.SeekStart); err == nil {
		t.Error("Seek(-1, SeekStart) succeeded, want error")
	}
}

func TestMemFSMkdirAllInvalid(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}
	if err := fsys.MkdirAll("invalid/../path", fs.ModePerm); err == nil {
		t.Error("MkdirAll(invalid/../path) succeeded, want error")
	}
}

func TestMemFSClose(t *testing.T) {
	fsys, err := newMemFS("memory://test")
	if err != nil {
		t.Fatal(err)
	}

	if err := fsys.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
