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
	"fmt"
	"strings"
	"testing"
)

func TestIsDirName(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input string
		want  bool
	}{
		{input: "", want: true},
		{input: "/", want: true},
		{input: "dir/", want: true},
		{input: "a/b/c/", want: true},
		{input: "file.txt", want: false},
		{input: "dir/file.txt", want: false},
		{input: ".", want: false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.input), func(t *testing.T) {
			t.Parallel()
			if got := isDirName(tc.input); got != tc.want {
				t.Errorf("isDirName(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidPath(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input   string
		wantErr bool
	}{
		{input: ".", wantErr: false},
		{input: "./.", wantErr: true},
		{input: "a\\b\\.\\..\\c", wantErr: false},
		{input: "a/b/./../c", wantErr: true},
		{input: "a/b/../c", wantErr: true},
		{input: "C:/", wantErr: true},
		{input: "C:\\", wantErr: false},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			err := validPath("open", tc.input)

			if tc.wantErr {
				if err == nil || !strings.Contains(err.Error(), "is not a valid path for") {
					t.Errorf("validPath(open, %s) expected to contain 'is not a valid path for', got: %q", tc.input, err)
				}
			} else {
				if err != nil {
					t.Errorf("validPath(open, %s) returned error %q, want nil", tc.input, err)
				}
			}
		})
	}
}
