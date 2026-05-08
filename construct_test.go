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

func TestNewBaseFS(t *testing.T) {
	tests := []struct {
		uri      string
		wantType string
		wantErr  bool
	}{
		{
			uri:      "angry://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "angryfs://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "file://",
			wantType: reflect.TypeFor[*localFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      ".",
			wantType: reflect.TypeFor[*localFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "mem://",
			wantType: reflect.TypeFor[*memFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "memfs://",
			wantType: reflect.TypeFor[*memFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "null://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
		{
			uri:      "nullfs://",
			wantType: reflect.TypeFor[*nullFS]().Name(),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
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
	}
}
