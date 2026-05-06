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
	"time"
)

var (
	_ FileInfo = (*fsInfo)(nil)
)

type fsInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (finfo *fsInfo) Name() string {
	return finfo.name
}

func (finfo *fsInfo) Size() int64 {
	return finfo.size
}

func (finfo *fsInfo) Mode() fs.FileMode {
	return finfo.mode
}

func (finfo *fsInfo) ModTime() time.Time {
	return finfo.modTime
}

func (finfo *fsInfo) IsDir() bool {
	return finfo.isDir
}

func (finfo *fsInfo) Sys() any {
	return finfo.sys
}
