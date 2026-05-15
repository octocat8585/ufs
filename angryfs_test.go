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

func TestIsAngryFSUri(t *testing.T) {
	testCases := []struct {
		name string
		want bool
	}{
		{name: "angry:", want: true},
		{name: "angry://", want: true},
		{name: "angryfs://", want: false},
		{name: "mem://", want: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isAngryFSUri(tc.name)
			if got != tc.want {
				t.Errorf("got: %t, want: %t", got, tc.want)
			}
		})
	}
}

func TestNewAngryFS(t *testing.T) {
	fsys := mustAngryFS(t)
	if fsys == nil {
		t.Fatal("fsys is nil")
	}
}

func TestAngryFSOpen(t *testing.T) {
	for _, tc := range testassetFilenameList {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			_, err := fsys.Open(tc)
			if err != fs.ErrInvalid {
				t.Errorf("Open(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSClose(t *testing.T) {
	fsys := mustAngryFS(t)
	if err := fsys.Close(); err != fs.ErrInvalid {
		t.Errorf("Close() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSMkdirAll(t *testing.T) {
	for _, tc := range []string{"a", "a/b", "a/b/c", "abc", "null"} {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if err := fsys.MkdirAll(tc, fs.ModePerm); err != fs.ErrInvalid {
				t.Errorf("MkdirAll(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSReadFile(t *testing.T) {
	for _, tc := range testassetFilenameList {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if _, err := fsys.ReadFile(tc); err != fs.ErrInvalid {
				t.Errorf("ReadFile(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSReadLink(t *testing.T) {
	for _, tc := range testassetFilenameList {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if _, err := fsys.ReadLink(tc); err != fs.ErrInvalid {
				t.Errorf("ReadLink(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSLstat(t *testing.T) {
	for _, tc := range testassetFilenameList {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if _, err := fsys.Lstat(tc); err != fs.ErrInvalid {
				t.Errorf("Lstat(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSReadDir(t *testing.T) {
	for tcDirName := range testassetDirList {
		t.Run(tcDirName, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if _, err := fsys.ReadDir(tcDirName); err != fs.ErrInvalid {
				t.Errorf("ReadDir(%q) = %v, want %v", tcDirName, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSGlob(t *testing.T) {
	fsys := mustAngryFS(t)
	if _, err := fsys.Glob("*.txt"); err != fs.ErrInvalid {
		t.Errorf("Glob() = %v, want %v", err, fs.ErrInvalid)
	}
}

func TestAngryFSCreate(t *testing.T) {
	for _, tc := range testassetCreateFileList {
		t.Run(tc, func(t *testing.T) {
			fsys := mustAngryFS(t)
			if _, err := fsys.Create(tc); err != fs.ErrInvalid {
				t.Errorf("Create(%q) = %v, want %v", tc, err, fs.ErrInvalid)
			}
		})
	}
}

func TestAngryFSString(t *testing.T) {
	fsys := mustAngryFS(t)
	if got := fsys.String(); got != angryFSPrefix {
		t.Errorf("String() got: %q, want %q", got, angryFSPrefix)
	}
}

func mustAngryFS(tb testing.TB) *angryFS {
	fsys, err := newAngryFS(angryFSPrefix)
	if err != nil {
		tb.Fatalf("newAngryFS() returned error, %s", err)
	}
	return fsys
}
