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
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ FS            = (*nestFS)(nil)
	_ fs.ReadDirFS  = (*nestFS)(nil)
	_ fs.ReadFileFS = (*nestFS)(nil)
	_ fs.ReadLinkFS = (*nestFS)(nil)
	_ fs.GlobFS     = (*nestFS)(nil)
	_ fs.StatFS     = (*nestFS)(nil)
)

func removePathComponent(name string, mountPath string) string {
	return strings.TrimLeft(strings.TrimPrefix(name, mountPath), "\\/")
}

func getPotentialArchives(name string) []string {
	components := strings.Split(name, string(os.PathSeparator))
	potentials := []string{}
	for idx, component := range components {
		if strings.HasSuffix(component, ".d") {
			potentials = append(potentials, filepath.Join(components[0:idx+1]...))
		}
	}
	return potentials
}

// nestFS is a wrapper for a base FS that supports automatic mounting of archives.
// This means that any archive can opened and read automatically. Archives are revealed via $filename.d name pattern.
type nestFS struct {
	fsys  FS
	fsMap map[string]*nestFS
}

func (fsys *nestFS) mountArchive(name string) (*nestFS, error) {
	ctx := context.Background()
	lfs, ok := fsys.fsys.(*localFS)
	var newFS *archiveFS
	if ok {
		absName, err := lfs.getAbsPath(name)
		if err != nil {
			return nil, &fs.PathError{
				Op:   "mount",
				Path: name,
				Err:  err,
			}
		}
		newFS, err = newArchiveFSFromLocalFS(ctx, absName)
		if err != nil {
			return nil, &fs.PathError{
				Op:   "mount",
				Path: name,
				Err:  err,
			}
		}
	} else {
		f, err := fsys.Open(name)
		if err != nil {
			return nil, &fs.PathError{
				Op:   "mount",
				Path: name,
				Err:  err,
			}
		}
		newFS, err = newArchiveFSFromFile(f)
		if err != nil {
			return nil, &fs.PathError{
				Op:   "mount",
				Path: name,
				Err:  err,
			}
		}
	}

	wrapped := createNestFS(newFS)
	fsys.fsMap[name+".d"] = wrapped
	return wrapped, nil
}

func (fsys *nestFS) getFSAndSubpath(name string) (*nestFS, string, error) {
	targetFS := fsys
	targetName := name
	for mountPath, subFS := range fsys.fsMap {
		subPath := removePathComponent(name, mountPath)
		if len(subPath) < len(targetName) {
			targetName = subPath
			targetFS = subFS
		}
	}

	archiveDirNames := getPotentialArchives(targetName)
	for _, archiveDirName := range archiveDirNames {
		archiveName := strings.TrimSuffix(archiveDirName, ".d")
		info, err := targetFS.Stat(archiveName)
		if info != nil && err == nil {
			subPath := removePathComponent(targetName, archiveDirName)
			subFS, err := targetFS.mountArchive(archiveName)
			if err != nil {
				return nil, "", &fs.PathError{
					Op:   "mount",
					Path: name,
					Err:  fmt.Errorf("cannot mount archive %s, %w", archiveName, err),
				}
			}
			targetFS, targetName, err := subFS.getFSAndSubpath(subPath)
			if err != nil {
				return nil, "", &fs.PathError{
					Op:   "mount",
					Path: name,
					Err:  fmt.Errorf("cannot mount archive %s, %w", archiveName, err),
				}
			}
			return targetFS, targetName, nil
		}
	}

	return targetFS, targetName, nil
}

func (fsys *nestFS) Open(name string) (fs.File, error) {
	if err := validPath("open", name); err != nil {
		return nil, err
	}

	mountFS, subName, err := fsys.getFSAndSubpath(name)
	if err != nil {
		return nil, err
	}

	return mountFS.fsys.Open(subName)
}

func (fsys *nestFS) Close() error {
	return fsys.fsys.Close()
}

func (fsys *nestFS) Create(name string) (File, error) {
	if err := validPath("create", name); err != nil {
		return nil, err
	}

	// TODO:
	return fsys.fsys.Create(name)
}

func (fsys *nestFS) MkdirAll(name string, perm fs.FileMode) error {
	// TODO:
	return fsys.fsys.MkdirAll(name, perm)
}

func (fsys *nestFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if cFsys, ok := fsys.fsys.(fs.ReadDirFS); ok {
		return cFsys.ReadDir(name)
	}

	f, err := fsys.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	readDirFile, ok := f.(fs.ReadDirFile)
	if !ok {
		return nil, &fs.PathError{
			Op:   "readDir",
			Path: name,
			Err:  fmt.Errorf("%s is not a directory", name),
		}
	}
	return readDirFile.ReadDir(-1)
}

func (fsys *nestFS) Stat(name string) (fs.FileInfo, error) {
	if cFsys, ok := fsys.fsys.(fs.StatFS); ok {
		return cFsys.Stat(name)
	}
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

func (fsys *nestFS) ReadFile(name string) ([]byte, error) {
	mountFS, subName, err := fsys.getFSAndSubpath(name)
	if err != nil {
		return nil, err
	}

	if cFsys, ok := mountFS.fsys.(fs.ReadFileFS); ok {
		return cFsys.ReadFile(subName)
	}
	return fs.ReadFile(mountFS.fsys, subName)
}

func (fsys *nestFS) ReadLink(name string) (string, error) {
	if cFsys, ok := fsys.fsys.(fs.ReadLinkFS); ok {
		return cFsys.ReadLink(name)
	}
	return "", &fs.PathError{Op: "readlink", Path: name, Err: fs.ErrInvalid}
}

func (fsys *nestFS) Lstat(name string) (fs.FileInfo, error) {
	if cFsys, ok := fsys.fsys.(fs.ReadLinkFS); ok {
		return cFsys.Lstat(name)
	}
	return fsys.Stat(name)
}

func (fsys *nestFS) Glob(pattern string) ([]string, error) {
	if cFsys, ok := fsys.fsys.(fs.GlobFS); ok {
		return cFsys.Glob(pattern)
	}
	return globFS(fsys, pattern)
}

func newNestFS(name string) (FS, error) {
	fsys, err := newBaseFS(name)
	if err != nil {
		return nil, err
	}
	return createNestFS(fsys), nil
}

func createNestFS(fsys FS) *nestFS {
	return &nestFS{
		fsys:  fsys,
		fsMap: map[string]*nestFS{},
	}
}
