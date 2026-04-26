//This file is part of xtagger. ©2023 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package filesystem

import (
	"context"
	"io/fs"
	"path/filepath"
	"os"
)

// Options for PushMetas
type PushOpts struct {
	FollowSymlinks bool //Follow symbolic links if true
	Recursive bool //Only read contents of root directory, skip subdirectories.
}

type Meta struct {
	dir string //File's parent directory
	fs.FileInfo
}

// Parent directory (absolute)
func (r *Meta) Dir() string {
	return r.dir
}

// Path to file (can be relative)
func (r *Meta) Path() string {
	return filepath.Join(r.dir, r.Name())
}

// Return true if regular file
func (r *Meta) IsRegular() bool {
	return r.Mode().IsRegular()
}

// Return true if symbolic link
func (r *Meta) IsLink() bool {
	if r.Mode() & fs.ModeSymlink == fs.ModeSymlink {
		return true
	}
	return false
}

// Read file metadata for the given path
func NewMeta(dir string, info fs.FileInfo) *Meta {
	return &Meta{
		FileInfo: info,
		dir: dir,
	}
}

func Lstat(path string) (*Meta, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	return NewMeta(filepath.Dir(path), info), nil
}

// PushMetas is a producer, it starts a goroutine that stats files and 
// pushes the corresponding *Meta objects to the returned channel.
// Channel errs exists for returning errors to the error consumer.
// Variable start is the entry path for the stat operation.
// Variable opts controls options like recursion into subdirs and
// symlink behaviour.
// IMPORTANT: The caller (consumer) is responsible for draining the
// returned channel to prevent a goroutine leak.
func PushMetas(ctx context.Context, errs chan<- error, opts PushOpts, root string) <-chan *Meta {
	push := make(chan *Meta)
	go pushMetasAndClose(ctx, opts, root, push, errs)
	return push
}

// Wraps call to pushMetas to ensure channel closure only once as
// pushMetas can be called recursively.
func pushMetasAndClose(ctx context.Context, opts PushOpts, root string, pusher chan<- *Meta, errs chan<- error) {
	defer close(pusher)
	pushMetas(ctx, opts, root, pusher, errs)
}

// Walks down path starting at root, stats every regular file, builds
// and pushes *Meta objects to the pusher channel.
func pushMetas(ctx context.Context, opts PushOpts, root string, pusher chan<- *Meta, errs chan<- error) {
	walker := func(path string, d fs.DirEntry, err error)error {
		// Check for cancellation
		select {
			case <-ctx.Done():
			return fs.SkipAll
			default:
		}
		// Handle directory entries
		if d.IsDir() {
			if !opts.Recursive && path != root {
				return fs.SkipDir
			}
			return nil //skip directory entries
		}
		// Handle symlinks
		if d.Type() & fs.ModeSymlink == fs.ModeSymlink {
			if !opts.FollowSymlinks {
				return nil
			}
			resolvedPath, err := resolveSymlink(filepath.Split(path))
			if err != nil {
				errs <- err
				return nil
			}
			pushMetas(ctx, opts, resolvedPath, pusher, errs)
		}
		// Handle regular files
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				errs <- err
				return nil
			}
			pusher <- NewMeta(filepath.Dir(path), info)
		}
		return nil
	}
	err := filepath.WalkDir(root, walker)
	if err != nil {
		errs <- err
	}
}

// Resolve the symlink name in directory dir.
// If the symlink is relative, the relative path will be appended to
// dir, then filepath.Clean() will be applied.
func resolveSymlink(dir, name string) (string, error) {
	resolved, err := os.Readlink(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(resolved) {
		return resolved, nil
	}
	return filepath.Clean(filepath.Join(dir, resolved)), nil
}
