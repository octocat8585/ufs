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

//go:build windows

package ufs

import (
	"os"
	"testing"
)

// TestNameToURIWindowsDriveLetters verifies that nameToURI correctly converts
// "file://X:\path" strings (which url.Parse rejects on Windows) into valid
// file:// URL structs for drive letters other than C:.
func TestNameToURIWindowsDriveLetters(t *testing.T) {
	tests := []struct {
		input    string
		wantPath string
	}{
		{`file://D:\`, `/D:/`},
		{`file://D:\some\path`, `/D:/some/path`},
		{`file://Z:\data`, `/Z:/data`},
		{`file://Z:\deeply\nested\dir`, `/Z:/deeply/nested/dir`},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			u, err := nameToURI(tt.input)
			if err != nil {
				t.Fatalf("nameToURI(%q) = %v, want nil", tt.input, err)
			}
			if u.Scheme != "file" {
				t.Errorf("scheme = %q, want %q", u.Scheme, "file")
			}
			if u.Path != tt.wantPath {
				t.Errorf("path = %q, want %q", u.Path, tt.wantPath)
			}
		})
	}
}

// TestLocalFSNormalizePathWindowsDriveLetters verifies that localFSNormalizePath
// converts the "/D:/path" form (produced after stripping "file://" from a
// canonical file:///D:/path URI) to the Windows-native "D:\path" form.
func TestLocalFSNormalizePathWindowsDriveLetters(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`file:///D:/some/path`, `D:\some\path`},
		{`file:///D:/`, `D:\`},
		{`file:///Z:/data`, `Z:\data`},
		{`file:///Z:/deeply/nested/dir`, `Z:\deeply\nested\dir`},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := localFSNormalizePath(tt.input)
			if got != tt.want {
				t.Errorf("localFSNormalizePath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestNewWindowsDriveLetters verifies that New() can mount drives other than C:\
// using both a bare Windows path ("D:\") and a file://-prefixed URI ("file://D:\").
// Sub-tests are skipped when the drive is not present on the host.
func TestNewWindowsDriveLetters(t *testing.T) {
	drives := []struct {
		bare    string
		fileURI string
	}{
		{`D:\`, `file://D:\`},
		{`Z:\`, `file://Z:\`},
	}
	for _, d := range drives {
		d := d
		t.Run(d.bare, func(t *testing.T) {
			if _, err := os.Stat(d.bare); err != nil {
				t.Skipf("drive %s not available: %v", d.bare, err)
			}
			for _, name := range []string{d.bare, d.fileURI} {
				fsys, err := New(name)
				if err != nil {
					t.Errorf("New(%q) = %v, want nil", name, err)
					continue
				}
				_ = fsys.Close()
			}
		})
	}
}
