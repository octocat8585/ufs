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
	"time"
)

func TestInfoDir(t *testing.T) {
	tNow := mustTime("2026-01-01T00:00:00Z")
	sysVal := struct{ x int }{x: 42}
	info := &fsInfo{
		name:    "mydir",
		size:    4096,
		mode:    fs.ModeDir | fs.ModePerm,
		modTime: tNow,
		isDir:   true,
		sys:     sysVal,
	}
	if info.Name() != "mydir" {
		t.Errorf("Name() = %q, want %q", info.Name(), "mydir")
	}
	if info.Size() != 4096 {
		t.Errorf("Size() = %d, want 4096", info.Size())
	}
	if info.Mode() != fs.ModeDir|fs.ModePerm {
		t.Errorf("Mode() = %v, want %v", info.Mode(), fs.ModeDir|fs.ModePerm)
	}
	if info.ModTime() != tNow {
		t.Errorf("ModTime() = %v, want %v", info.ModTime(), tNow)
	}
	if !info.IsDir() {
		t.Error("IsDir() = false, want true")
	}
	if info.Sys() != sysVal {
		t.Errorf("Sys() = %v, want %v", info.Sys(), sysVal)
	}
}

func TestInfo(t *testing.T) {
	tNow := time.Now()
	info := &fsInfo{
		name:    "test",
		size:    0,
		mode:    0,
		modTime: tNow,
		isDir:   false,
		sys:     nil,
	}

	if info.Name() != "test" {
		t.Errorf("Name() = %q, want %q", info.Name(), "test")
	}
	if info.Size() != 0 {
		t.Errorf("Size() = %d, want %d", info.Size(), 0)
	}
	if info.Mode() != 0 {
		t.Errorf("Mode() = %d, want %d", info.Mode(), 0)
	}
	if info.ModTime() != tNow {
		t.Errorf("ModTime() = %v, want %v", info.ModTime(), tNow)
	}
	if info.IsDir() != false {
		t.Errorf("IsDir() = %v, want %v", info.IsDir(), false)
	}
	if info.Sys() != nil {
		t.Errorf("Sys() = %v, want %v", info.Sys(), nil)
	}
}
