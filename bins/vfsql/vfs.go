// filepath: /home/lxk/Desktop/cirm/bins/vfsql/vfs.go
package vfsql

import (
	"database/sql"
	"errors"
	"os"
	"strings"
)

// Open opens an existing file for reading
func (vfs *VFS) Open(path string) (*File, error) {
	return vfs.OpenFile(path, os.O_RDONLY, 0)
}

// OpenFile opens a file with specified flags and permissions
func (vfs *VFS) OpenFile(path string, flag int, perm os.FileMode) (*File, error) {
	path = normalizePath(path)

	// Check for variant path
	if isVariantPath(path) {
		origPath, variantName, ok := splitVariantPath(path)
		if !ok {
			return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
		}
		return vfs.GetVariant(origPath, variantName)
	}

	// Resolve the inode
	inode, err := vfs.resolvePath(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && (flag&os.O_CREATE) != 0 {
			return vfs.createFile(path, perm, flag)
		}
		return nil, err
	}

	// Check if it's a directory being opened as file
	if inode.typ == "dir" && (flag&(os.O_WRONLY|os.O_RDWR)) != 0 {
		return nil, &os.PathError{Op: "open", Path: path, Err: errors.New("is a directory")}
	}

	// Truncate if requested
	if (flag & os.O_TRUNC) != 0 {
		if err := vfs.truncateFile(inode); err != nil {
			return nil, err
		}
	}

	// Load content if it's a file
	var content []byte
	if inode.typ == "file" {
		err := vfs.db.db.QueryRow(
			"SELECT content FROM storage WHERE inode_id = ?",
			inode.id,
		).Scan(&content)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	writable := (flag & (os.O_WRONLY | os.O_RDWR)) != 0

	file := &File{
		vfs:      vfs,
		inode:    inode,
		path:     path,
		content:  content,
		pos:      0,
		flags:    flag,
		writable: writable,
		dirty:    false,
		closed:   false,
	}

	// If append mode, seek to end
	if (flag & os.O_APPEND) != 0 {
		file.pos = int64(len(content))
	}

	return file, nil
}

// Create creates a new file
func (vfs *VFS) Create(path string) (*File, error) {
	return vfs.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

// createFile is a helper to create a new file
func (vfs *VFS) createFile(path string, perm os.FileMode, flag int) (*File, error) {
	dir, name := splitPath(path)

	// Ensure parent directory exists
	parentInode, err := vfs.resolvePath(dir)
	if err != nil {
		return nil, &os.PathError{Op: "create", Path: path, Err: err}
	}

	if parentInode.typ != "dir" {
		return nil, &os.PathError{Op: "create", Path: path, Err: errors.New("parent is not a directory")}
	}

	tx, err := vfs.db.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Create inode
	result, err := tx.Exec(
		"INSERT INTO inodes (fs_id, parent_id, name, type, mode, uid, gid, size, created_at, modified_at, accessed_at) VALUES (?, ?, ?, 'file', ?, ?, ?, 0, ?, ?, ?)",
		vfs.id, parentInode.id, name, perm&os.ModePerm, 0, 0, now, now, now,
	)
	if err != nil {
		return nil, &os.PathError{Op: "create", Path: path, Err: err}
	}

	inodeID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create empty storage entry
	_, err = tx.Exec(
		"INSERT INTO storage (inode_id, content) VALUES (?, ?)",
		inodeID, []byte{},
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	inode := &inode{
		id:         inodeID,
		fsID:       vfs.id,
		parentID:   &parentInode.id,
		name:       name,
		typ:        "file",
		mode:       perm & os.ModePerm,
		uid:        0,
		gid:        0,
		size:       0,
		createdAt:  now,
		modifiedAt: now,
		accessedAt: now,
	}

	// Emit create event
	vfs.emitEvent(Event{
		Type:      EventCreate,
		Path:      path,
		Timestamp: now,
		FileType:  FileTypeFile,
		Size:      0,
	})

	writable := (flag & (os.O_WRONLY | os.O_RDWR)) != 0

	return &File{
		vfs:      vfs,
		inode:    inode,
		path:     path,
		content:  []byte{},
		pos:      0,
		flags:    flag,
		writable: writable,
		dirty:    false,
		closed:   false,
	}, nil
}

// truncateFile truncates a file to zero size
func (vfs *VFS) truncateFile(inode *inode) error {
	tx, err := vfs.db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Update inode
	_, err = tx.Exec(
		"UPDATE inodes SET size = 0, modified_at = ? WHERE id = ?",
		now, inode.id,
	)
	if err != nil {
		return err
	}

	// Update storage
	_, err = tx.Exec(
		"UPDATE storage SET content = ? WHERE inode_id = ?",
		[]byte{}, inode.id,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	inode.size = 0
	inode.modifiedAt = now

	return nil
}

// Remove removes a file or empty directory
func (vfs *VFS) Remove(path string) error {
	path = normalizePath(path)

	if path == "/" {
		return &os.PathError{Op: "remove", Path: path, Err: errors.New("cannot remove root directory")}
	}

	inode, err := vfs.resolvePath(path)
	if err != nil {
		return &os.PathError{Op: "remove", Path: path, Err: err}
	}

	// Check if directory is empty
	if inode.typ == "dir" {
		var count int
		err := vfs.db.db.QueryRow(
			"SELECT COUNT(*) FROM inodes WHERE parent_id = ?",
			inode.id,
		).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			return &os.PathError{Op: "remove", Path: path, Err: errors.New("directory not empty")}
		}
	}

	// Delete the inode (cascade will handle storage and variants)
	_, err = vfs.db.db.Exec("DELETE FROM inodes WHERE id = ?", inode.id)
	if err != nil {
		return err
	}

	// Emit delete event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventDelete,
		Path:      path,
		Timestamp: currentTimestamp(),
		FileType:  fileType,
	})

	return nil
}

// RemoveAll removes a path and all its children
func (vfs *VFS) RemoveAll(path string) error {
	path = normalizePath(path)

	if path == "/" {
		return &os.PathError{Op: "removeall", Path: path, Err: errors.New("cannot remove root directory")}
	}

	inode, err := vfs.resolvePath(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // Already doesn't exist
		}
		return &os.PathError{Op: "removeall", Path: path, Err: err}
	}

	// Delete the inode (cascade will handle everything)
	_, err = vfs.db.db.Exec("DELETE FROM inodes WHERE id = ?", inode.id)
	if err != nil {
		return err
	}

	// Emit delete event
	fileType := FileTypeFile
	if inode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventDelete,
		Path:      path,
		Timestamp: currentTimestamp(),
		FileType:  fileType,
	})

	return nil
}

// Rename renames (moves) a file or directory
func (vfs *VFS) Rename(oldpath, newpath string) error {
	oldpath = normalizePath(oldpath)
	newpath = normalizePath(newpath)

	if oldpath == "/" || newpath == "/" {
		return &os.PathError{Op: "rename", Path: oldpath, Err: errors.New("cannot rename root directory")}
	}

	// Get source inode
	srcInode, err := vfs.resolvePath(oldpath)
	if err != nil {
		return &os.PathError{Op: "rename", Path: oldpath, Err: err}
	}

	// Check if destination exists
	_, err = vfs.resolvePath(newpath)
	if err == nil {
		return &os.PathError{Op: "rename", Path: newpath, Err: os.ErrExist}
	}

	// Get new parent directory
	newDir, newName := splitPath(newpath)
	newParentInode, err := vfs.resolvePath(newDir)
	if err != nil {
		return &os.PathError{Op: "rename", Path: newpath, Err: err}
	}

	if newParentInode.typ != "dir" {
		return &os.PathError{Op: "rename", Path: newpath, Err: errors.New("parent is not a directory")}
	}

	// Update inode
	now := currentTimestamp()
	_, err = vfs.db.db.Exec(
		"UPDATE inodes SET parent_id = ?, name = ?, modified_at = ? WHERE id = ?",
		newParentInode.id, newName, now, srcInode.id,
	)
	if err != nil {
		return err
	}

	// Emit rename event
	fileType := FileTypeFile
	if srcInode.typ == "dir" {
		fileType = FileTypeDir
	}

	vfs.emitEvent(Event{
		Type:      EventRename,
		Path:      newpath,
		OldPath:   oldpath,
		Timestamp: now,
		FileType:  fileType,
	})

	return nil
}

// resolvePath resolves a path to an inode
func (vfs *VFS) resolvePath(path string) (*inode, error) {
	path = normalizePath(path)

	// Handle root
	if path == "/" {
		return vfs.getRootInode()
	}

	// Split path into components
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Start from root
	currentID := vfs.rootInodeID

	// Traverse path
	for _, part := range parts {
		if part == "" {
			continue
		}

		var id int64
		var typ string
		err := vfs.db.db.QueryRow(
			"SELECT id, type FROM inodes WHERE fs_id = ? AND parent_id = ? AND name = ?",
			vfs.id, currentID, part,
		).Scan(&id, &typ)

		if err == sql.ErrNoRows {
			return nil, os.ErrNotExist
		}
		if err != nil {
			return nil, err
		}

		currentID = id
	}

	// Load full inode data
	return vfs.getInodeByID(currentID)
}

// getRootInode returns the root inode
func (vfs *VFS) getRootInode() (*inode, error) {
	return vfs.getInodeByID(vfs.rootInodeID)
}

// getInodeByID loads an inode by ID
func (vfs *VFS) getInodeByID(id int64) (*inode, error) {
	var inode inode
	var parentID sql.NullInt64

	err := vfs.db.db.QueryRow(
		"SELECT id, fs_id, parent_id, name, type, mode, uid, gid, size, COALESCE(description, ''), COALESCE(tags, ''), created_at, modified_at, accessed_at FROM inodes WHERE id = ?",
		id,
	).Scan(
		&inode.id, &inode.fsID, &parentID, &inode.name, &inode.typ,
		&inode.mode, &inode.uid, &inode.gid, &inode.size,
		&inode.description, &inode.tags,
		&inode.createdAt, &inode.modifiedAt, &inode.accessedAt,
	)

	if err == sql.ErrNoRows {
		return nil, os.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		inode.parentID = &parentID.Int64
	}

	return &inode, nil
}
