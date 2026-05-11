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
	"os"
	"runtime"
	"strings"
)

const (
	pathSeparator     = string(os.PathSeparator)
	unixPathSeparator = "/"
)

func trimSlash(name string) string {
	return strings.Trim(name, "\\/")
}

func splitPath(name string) []string {
	return strings.Split(trimSlash(name), pathSeparator)
}

func validPath(op string, name string) error {
	if !fs.ValidPath(name) {
		return pathError(op, name, fmt.Errorf("%q is not a valid path for %s, %w", name, runtime.GOOS, fs.ErrInvalid))
	}
	return nil
}

func coerceUnix(name string) string {
	return strings.ReplaceAll(name, "\\", unixPathSeparator)
}

func isDirName(name string) bool {
	return isCwd(name) || strings.HasSuffix(name, pathSeparator)
}

func isCwd(name string) bool {
	return name == "" || name == "."
}

func pathError(op string, name string, err error) error {
	return &fs.PathError{
		Op:   op,
		Path: name,
		Err:  err,
	}
}
