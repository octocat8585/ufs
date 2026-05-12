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
	"io/fs"
	"reflect"
	"testing"
)

func TestNewBaseFSInvalid(t *testing.T) {
	_, err := newBaseFS("unknown://surely-not-a-path-xyz")
	if err == nil {
		t.Error("newBaseFS(unknown://) succeeded, want error")
	}
	if _, ok := err.(*fs.PathError); !ok {
		t.Errorf("newBaseFS(unknown://) returned %T, want *fs.PathError", err)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		uri      string
		wantType string
		wantErr  bool
		nested   bool
	}{
		{
			uri:      "angry://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "angry://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "file://",
			wantType: reflect.TypeFor[*localFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      cwdPath,
			wantType: reflect.TypeFor[*localFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "memory://",
			wantType: reflect.TypeFor[*memFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "memory://",
			wantType: reflect.TypeFor[*memFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "null://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "null://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "file:///?a=memory://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
			nested:   true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("newBaseFS(%q)", tt.uri), func(t *testing.T) {
			if tt.nested {
				t.Skip("test case requires nestFS support")
			}
			got, err := newBaseFS(tt.uri)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getBaseFS(%q) = %q, want error", tt.uri, got)
				}
			} else {
				if err != nil {
					t.Errorf("getBaseFS(%q) = %q, want %q", tt.uri, err, tt.wantType)
				} else {
					gotTypeName := reflect.TypeOf(got).Name()
					if gotTypeName != tt.wantType {
						t.Errorf("getBaseFS(%q) = %q, want %q", tt.uri, got, tt.wantType)
					}
				}
			}
		})
		t.Run(fmt.Sprintf("New(%q)", tt.uri), func(t *testing.T) {
			got, err := New(tt.uri)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getBaseFS(%q) = %q, want error", tt.uri, got)
				}
			} else {
				if err != nil {
					t.Errorf("getBaseFS(%q) = %q, want %q", tt.uri, err, tt.wantType)
				} else {
					gotAsNestFS, ok := got.(*nestFS)
					if ok {
						gotTypeName := reflect.TypeOf(gotAsNestFS.fsys).Name()
						if gotTypeName != tt.wantType {
							t.Errorf("getBaseFS(%q) = %q, want %q", tt.uri, got, tt.wantType)
						}
					} else {
						t.Errorf("%q is not of type *nestFS", got)
					}
				}
			}
		})
	}
}

func TestCreateURI(t *testing.T) {
	tests := []struct {
		name   string
		mounts map[string]string
		want   string
	}{
		{
			name:   "memory://",
			mounts: map[string]string{},
			want:   "memory:",
		},
		{
			name:   "angry://",
			mounts: map[string]string{},
			want:   "angry:",
		},
		{
			name:   "null://",
			mounts: map[string]string{},
			want:   "null:",
		},
		{
			name:   "file://",
			mounts: map[string]string{},
			want:   "file:",
		},
		{
			name: "file:///",
			mounts: map[string]string{
				"mounted/null":   "null://",
				"mounted/angry":  "angry://",
				"mounted/memory": "memory://",
			},
			want: "file:///?mounted%2Fangry=angry%3A&mounted%2Fmemory=memory%3A&mounted%2Fnull=null%3A",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CreateURI(tc.name, tc.mounts)
			if err != nil {
				t.Errorf("CreateURI(%q, %v) got error: %s", tc.name, tc.mounts, err)
			}
			if got != tc.want {
				t.Errorf("CreateURI(%q, %v) got: %q, want: %q", tc.name, tc.mounts, got, tc.want)
			}
			fsys, err := New(got)
			if err != nil {
				t.Fatalf("New(%q) got error: %s", got, err)
			}

			gotTypeName := reflect.TypeOf(fsys).Name()
			wantType := reflect.TypeFor[*nestFS]().Name()
			if gotTypeName != wantType {
				t.Errorf("New(%q) = %q, want %q", got, fsys, wantType)
			}
		})
	}
}
