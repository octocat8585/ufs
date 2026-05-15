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
	"io"
	"io/fs"
	"path/filepath"
	"sort"
)

// Rsync copies a directory structure of the source file system to the destination file system.
func Rsync(srcFS fs.FS, destFS FS, dir string) error {
	return ForEachFilename(srcFS, dir, func(name string) error {
		dir, _ := filepath.Split(name)
		dir = filepath.Clean(dir)
		if err := destFS.MkdirAll(dir, fs.ModePerm); err != nil {
			return err
		}
		if err := Copy(srcFS, name, destFS, name); err != nil {
			return err
		}
		return nil
	})
}

// Copy a file from one file system to another.
func Copy(srcFS fs.FS, srcFilename string, destFS FS, destFilename string) error {
	sfp, err := srcFS.Open(srcFilename)
	if err != nil {
		return err
	}
	defer sfp.Close()

	dfp, err := destFS.Create(destFilename)
	if err != nil {
		return err
	}
	defer dfp.Close()

	if _, err := io.Copy(dfp, sfp); err != nil {
		return err
	}
	return nil
}

// ForEachFilename is similar to ListFiles but calls against an object.
// Using this over ListFiles will save memory for large file listings.
func ForEachFilename(fsys fs.FS, dir string, f func(string) error) error {
	lf, ok := fsys.(ForEachFilenameIter)
	if ok {
		return lf.ForEachFilename(dir, f)
	}
	files, err := ListFiles(fsys, dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := f(file); err != nil {
			return err
		}
	}
	return nil
}

// ForEachFileInfo is similar to ListFiles but calls against an object.
// Using this over ListFiles will save memory for large file listings.
func ForEachFileInfo(fsys fs.FS, dir string, f func(fs.FileInfo) error) error {
	lf, ok := fsys.(ForEachFileInfoIter)
	if ok {
		return lf.ForEachFileInfo(dir, f)
	}
	files, err := ListFiles(fsys, dir)
	if err != nil {
		return err
	}
	for _, filename := range files {
		fileInfo, err := fs.Stat(fsys, filename)
		if err != nil {
			return err
		}
		if err := f(fileInfo); err != nil {
			return err
		}
	}
	return nil
}

// List directories and files.
func List(fsys fs.FS, dir string) ([]string, error) {
	return list(fsys, dir, true)
}

// ListFiles is similar to list but excludes directories.
func ListFiles(fsys fs.FS, dir string) ([]string, error) {
	lf, ok := fsys.(ListFilenames)
	if ok {
		return lf.ListFilenames(dir)
	}
	return list(fsys, dir, false)
}

func list(fsys fs.FS, dir string, includeDirs bool) ([]string, error) {
	m := map[string]any{}
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if isCwd(path) && (d == nil || err != nil) {
			return nil
		}
		if err != nil {
			return err
		}
		if includeDirs || !d.IsDir() {
			m[path] = nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	items := make([]string, len(m))
	i := 0
	for k := range m {
		items[i] = k
		i++
	}
	sort.Strings(items)
	return items, nil
}
