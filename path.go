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
	"path"
	"runtime"
	"strings"
)

const (
	unixPathSeparator         = "/"
	windowsPathSeparator      = "\\"
	cwdPath                   = "."
	unixAndWindowsSlashCutset = unixPathSeparator + windowsPathSeparator
	emptyDirSize              = 0
)

func removePathPrefix(name string, removePath string) (string, bool) {
	removePath = path.Clean(removePath)
	name = path.Clean(name)
	if isCwd(removePath) {
		return name, true
	}
	if removePath == name {
		return cwdPath, true
	}
	return strings.CutPrefix(name, removePath+unixPathSeparator)
}

func trimSlash(name string) string {
	return strings.Trim(name, unixAndWindowsSlashCutset)
}

func splitPath(name string) []string {
	return strings.Split(trimSlash(name), unixPathSeparator)
}

func validPath(op string, name string) error {
	if !fs.ValidPath(name) {
		return pathError(op, name, fmt.Errorf("%q is not a valid path for %s, %w", name, runtime.GOOS, fs.ErrInvalid))
	}
	return nil
}

func coerceUnix(name string) string {
	return strings.ReplaceAll(name, windowsPathSeparator, unixPathSeparator)
}

func isDirName(name string) bool {
	return isCwd(name) || strings.HasSuffix(name, unixPathSeparator)
}

func isCwd(name string) bool {
	return name == "" || name == cwdPath
}

func pathError(op string, name string, err error) error {
	return &fs.PathError{
		Op:   op,
		Path: name,
		Err:  err,
	}
}
