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
	"os"
	"strings"
	"testing"
)

func TestCreateOSTempDirectory(t *testing.T) {
	dir, cleanup, err := createOSTempDirectory()
	if err != nil {
		t.Error(err)
	}
	if !osExists(dir) {
		t.Errorf("'%s' does not exist when it should", dir)
	}

	if !strings.Contains(dir, "goapp") {
		t.Errorf("'%s' does not contain 'goapp'", dir)
	}
	cleanup()
	if osExists(dir) {
		t.Errorf("'%s' exists when it should not", dir)
	}
}

func TestOSDeleteFile(t *testing.T) {
	t.Run("nonexistent", func(t *testing.T) {
		err := osDeleteFile("/nonexistent/path/that/cannot/exist-" + t.Name() + ".txt")
		if err != nil {
			t.Errorf("osDeleteFile(nonexistent) = %v, want nil", err)
		}
	})

	t.Run("existing", func(t *testing.T) {
		f, err := os.CreateTemp("", "ufs-osutil-test-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		p := f.Name()
		f.Close()

		if err := osDeleteFile(p); err != nil {
			t.Errorf("osDeleteFile(existing) = %v, want nil", err)
		}
		if osExists(p) {
			t.Errorf("%q still exists after osDeleteFile", p)
		}
	})
}

func TestTryOSDeleteFile(t *testing.T) {
	f, err := os.CreateTemp("", "ufs-try-delete-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	p := f.Name()
	f.Close()

	tryOSDeleteFile(p)
	if osExists(p) {
		t.Errorf("%q still exists after tryOSDeleteFile", p)
	}
}

func TestOSDeleteDirectoryExists(t *testing.T) {
	dir, err := os.MkdirTemp("", "ufs-del-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	if err := osDeleteDirectory(dir); err != nil {
		t.Errorf("osDeleteDirectory(existing) = %v, want nil", err)
	}
	if osExists(dir) {
		t.Errorf("%q still exists after osDeleteDirectory", dir)
	}
}

func TestTryOSDeleteDirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "ufs-try-del-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	tryOSDeleteDirectory(dir)
	if osExists(dir) {
		t.Errorf("%q still exists after tryOSDeleteDirectory", dir)
	}
}

func TestNewRemoteArchive(t *testing.T) {
	fsys, err := New("https://github.com/mholt/archives/archive/refs/heads/main.zip")
	if err != nil {
		t.Error(err)
	}
	defer fsys.Close()

	if files, err := fsys.ReadDir(cwdPath); files != nil {
		t.Logf("files: %v, err: %s", files, err)
	}
	if files, err := fsys.ReadDir("archives-main"); files != nil {
		t.Logf("files: %v, err: %s", files, err)
	}
}
