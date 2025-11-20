// filepath: /home/lxk/Desktop/cirm/bins/vfsql/vfs_ops.go
package vfsql

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"time"
)

// Mkdir creates a directory
func (vfs *VFS) Mkdir(path string, perm os.FileMode) error {
	path = normalizePath(path)

	// Check if already exists
	_, err := vfs.resolvePath(path)
	if err == nil {
		return &os.PathError{Op: "mkdir", Path: path, Err: os.ErrExist}
	}

	// Get parent directory
	dir, name := splitPath(path)
	parentInode, err := vfs.resolvePath(dir)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: path, Err: err}
	}

	if parentInode.typ != "dir" {
		return &os.PathError{Op: "mkdir", Path: path, Err: errors.New("parent is not a directory")}
	}

	now := currentTimestamp()

	// Create directory inode
	_, err = vfs.db.db.Exec(
		"INSERT INTO inodes (fs_id, parent_id, name, type, mode, uid, gid, size, created_at, modified_at, accessed_at) VALUES (?, ?, ?, 'dir', ?, ?, ?, 0, ?, ?, ?)",
		vfs.id, parentInode.id, name, perm&os.ModePerm, 0, 0, now, now, now,
	)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: path, Err: err}
	}

	// Emit create event
	vfs.emitEvent(Event{
		Type:      EventCreate,
		Path:      path,
		Timestamp: now,
		FileType:  FileTypeDir,
	})

	return nil
}

// MkdirAll creates a directory and all necessary parent directories
func (vfs *VFS) MkdirAll(path string, perm os.FileMode) error {
	path = normalizePath(path)

	// Check if already exists
	inode, err := vfs.resolvePath(path)
	if err == nil {
		if inode.typ == "dir" {
			return nil // Already exists as directory
		}
		return &os.PathError{Op: "mkdir", Path: path, Err: errors.New("exists but not a directory")}
	}

	// Create parent directories recursively
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentPath := "/"

	for _, part := range parts {
		if part == "" {
			continue
		}

		currentPath = joinPath(currentPath, part)

		_, err := vfs.resolvePath(currentPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if err := vfs.Mkdir(currentPath, perm); err != nil && !errors.Is(err, os.ErrExist) {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

// ReadDir reads directory entries
func (vfs *VFS) ReadDir(path string) ([]fs.DirEntry, error) {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return nil, &os.PathError{Op: "readdir", Path: path, Err: err}
	}

	if inode.typ != "dir" {
		return nil, &os.PathError{Op: "readdir", Path: path, Err: errors.New("not a directory")}
	}

	rows, err := vfs.db.db.Query(
		"SELECT id, name, type, mode, size, modified_at FROM inodes WHERE fs_id = ? AND parent_id = ? ORDER BY name",
		vfs.id, inode.id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []fs.DirEntry
	for rows.Next() {
		var id int64
		var name, typ string
		var mode int64
		var size, modTime int64

		if err := rows.Scan(&id, &name, &typ, &mode, &size, &modTime); err != nil {
			return nil, err
		}

		info := &fileInfo{
			name:    name,
			size:    size,
			mode:    os.FileMode(mode),
			modTime: time.Unix(modTime, 0),
			isDir:   typ == "dir",
		}

		entries = append(entries, &dirEntry{info: info})
	}

	return entries, rows.Err()
}

// Stat returns file info
func (vfs *VFS) Stat(path string) (fs.FileInfo, error) {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: path, Err: err}
	}

	return &fileInfo{
		name:    inode.name,
		size:    inode.size,
		mode:    inode.mode,
		modTime: time.Unix(inode.modifiedAt, 0),
		isDir:   inode.typ == "dir",
		sys:     inode,
	}, nil
}

// Lstat returns file info (same as Stat for this VFS)
func (vfs *VFS) Lstat(path string) (fs.FileInfo, error) {
	return vfs.Stat(path)
}

// Chmod changes file permissions
func (vfs *VFS) Chmod(path string, mode os.FileMode) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "chmod", Path: path, Err: err}
	}

	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET mode = ?, modified_at = ? WHERE id = ?",
		mode&os.ModePerm, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit chmod event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventChmod,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
	})

	return nil
}

// Chown changes file ownership
func (vfs *VFS) Chown(path string, uid, gid int) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "chown", Path: path, Err: err}
	}

	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET uid = ?, gid = ?, modified_at = ? WHERE id = ?",
		uid, gid, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit chown event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventChown,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
	})

	return nil
}

// Getwd returns the current working directory
func (vfs *VFS) Getwd() (string, error) {
	return vfs.cwd, nil
}

// Chdir changes the current working directory
func (vfs *VFS) Chdir(path string) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "chdir", Path: path, Err: err}
	}

	if inode.typ != "dir" {
		return &os.PathError{Op: "chdir", Path: path, Err: errors.New("not a directory")}
	}

	vfs.cwd = path
	return nil
}

// SetDescription sets the description metadata
func (vfs *VFS) SetDescription(path string, description string) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "setdescription", Path: path, Err: err}
	}

	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET description = ?, modified_at = ? WHERE id = ?",
		description, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit metadata event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventMetadata,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
	})

	return nil
}

// GetDescription gets the description metadata
func (vfs *VFS) GetDescription(path string) (string, error) {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return "", &os.PathError{Op: "getdescription", Path: path, Err: err}
	}

	return inode.description, nil
}

// SetTags sets the tags metadata
func (vfs *VFS) SetTags(path string, tags []string) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "settags", Path: path, Err: err}
	}

	tagsStr := formatTagsString(tags)
	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET tags = ?, modified_at = ? WHERE id = ?",
		tagsStr, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit metadata event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventMetadata,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
		Tags:      tags,
	})

	return nil
}

// GetTags gets the tags metadata
func (vfs *VFS) GetTags(path string) ([]string, error) {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return nil, &os.PathError{Op: "gettags", Path: path, Err: err}
	}

	return parseTagsString(inode.tags), nil
}

// AddTag adds a single tag
func (vfs *VFS) AddTag(path string, tag string) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "addtag", Path: path, Err: err}
	}

	newTags := addTag(inode.tags, tag)
	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET tags = ?, modified_at = ? WHERE id = ?",
		newTags, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit metadata event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventMetadata,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
		Tags:      parseTagsString(newTags),
	})

	return nil
}

// RemoveTag removes a single tag
func (vfs *VFS) RemoveTag(path string, tag string) error {
	path = normalizePath(path)

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "removetag", Path: path, Err: err}
	}

	newTags := removeTag(inode.tags, tag)
	now := currentTimestamp()

	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET tags = ?, modified_at = ? WHERE id = ?",
		newTags, now, inode.id,
	)
	if err != nil {
		return err
	}

	// Emit metadata event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventMetadata,
		Path:      path,
		Timestamp: now,
		FileType:  fileType,
		Tags:      parseTagsString(newTags),
	})

	return nil
}
