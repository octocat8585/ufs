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

func newBaseFS(name string) (FS, error) {
	if strings.HasPrefix(name, "memfs://") || strings.HasPrefix(name, "mem://") {
		return newMemFS(name)
	}
	if strings.HasPrefix(name, "angryfs://") || strings.HasPrefix(name, "angry://") {
		return newAngryFS(name)
	}
	if strings.HasPrefix(name, "nullfs://") || strings.HasPrefix(name, "null://") {
		return newNullFS(name)
	}
	if strings.HasPrefix(name, "file://") {
		return newLocalFS(name)
	}
	if strings.HasPrefix(name, "gs://") {
		return newGCSFS(name)
	}

	stat, err := os.Stat(name)
	if err == nil && stat != nil {
		return newLocalFS(name)
	}
	return nil, &fs.PathError{
		Op:   "mount",
		Path: name,
		Err:  fmt.Errorf("%q is not a valid mount path for %s, %w", name, runtime.GOOS, err),
	}
}
