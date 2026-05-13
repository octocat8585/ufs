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

func xTestNewRemoteArchive(t *testing.T) {
	fsys, err := New("https://github.com/mholt/archives/archive/refs/heads/main.zip")
	defer fsys.Close()

	if err != nil {
		t.Error(err)
	}
	if files, err := fsys.ReadDir("."); files != nil {
		t.Logf("files: %v, err: %s", files, err)
	}
	if files, err := fsys.ReadDir("archives-main"); files != nil {
		t.Logf("files: %v, err: %s", files, err)
	}
	t.Fail()
}
