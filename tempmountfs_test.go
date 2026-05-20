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
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTempMountFS(t *testing.T) {
	fsys, err := newTempMountFS("test://", func(string) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if fsys == nil {
		t.Fatal("fsys is nil")
	}
	if err := fsys.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}

func TestTempMountFSFileSystem(t *testing.T) {
	testFileSystem(t, func(name string) (FS, error) {
		return newTempMountFS(name, func(string) error { return nil })
	}, "temp://")
}

func TestTempMountFSCleanup(t *testing.T) {
	var capturedDir string
	fsys, err := newTempMountFS("test://", func(dir string) error {
		capturedDir = dir
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !osExists(capturedDir) {
		t.Error("temp dir should exist before Close")
	}
	if err := fsys.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
	if osExists(capturedDir) {
		t.Error("temp dir should be deleted after Close")
	}
}

func TestTempMountFSCloseError(t *testing.T) {
	// Use an angry FS so that lfs.Close() returns an error.
	angry := makeAngryFS(angryFSPrefix)
	tfs := makeTempMountFS(angry, "test://", func() error { return nil })
	err := tfs.Close()
	if err == nil {
		t.Fatal("Close() = nil, want error from angry lfs")
	}
}

func TestTempMountFSPrepareError(t *testing.T) {
	var capturedDir string
	wantErr := errors.New("prepare failed")
	_, err := newTempMountFS("test://", func(dir string) error {
		capturedDir = dir
		return wantErr
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if capturedDir != "" && osExists(capturedDir) {
		t.Error("temp dir should be cleaned up after prepare error")
	}
}

func TestTempMountFSReadLink(t *testing.T) {
	skipTestOnWindows(t)
	var tempDir string
	fsys, err := newTempMountFS("test://", func(dir string) error {
		tempDir = dir
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	f, err := fsys.Create("target.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()

	if err := os.Symlink("target.txt", filepath.Join(tempDir, "link.txt")); err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	rlfs, ok := fsys.(fs.ReadLinkFS)
	if !ok {
		t.Fatal("tempMountFS does not implement fs.ReadLinkFS")
	}
	got, err := rlfs.ReadLink("link.txt")
	if err != nil {
		t.Fatalf("ReadLink failed: %v", err)
	}
	if got != "target.txt" {
		t.Errorf("ReadLink = %q, want %q", got, "target.txt")
	}
}

func TestTempMountFSLstat(t *testing.T) {
	skipTestOnWindows(t)
	var tempDir string
	fsys, err := newTempMountFS("test://", func(dir string) error {
		tempDir = dir
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	f, err := fsys.Create("lstat_file.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()

	if err := os.Symlink("lstat_file.txt", filepath.Join(tempDir, "lstat_link.txt")); err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	rlfs, ok := fsys.(fs.ReadLinkFS)
	if !ok {
		t.Fatal("tempMountFS does not implement fs.ReadLinkFS")
	}

	t.Run("regular_file", func(t *testing.T) {
		info, err := rlfs.Lstat("lstat_file.txt")
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
		info, err := rlfs.Lstat("lstat_link.txt")
		if err != nil {
			t.Fatalf("Lstat on symlink failed: %v", err)
		}
		if info.Mode()&fs.ModeSymlink == 0 {
			t.Errorf("Mode() missing ModeSymlink for symlink: %v", info.Mode())
		}
	})
}
