// filepath: /home/lxk/Desktop/cirm/bins/vfsql/vfs_variants.go
package vfsql

import (
	"database/sql"
	"errors"
	"os"
)

// CreateVariant creates a new variant of a file
func (vfs *VFS) CreateVariant(originalPath string, variantName string, content []byte) error {
	originalPath = normalizePath(originalPath)

	// Get original file inode
	origInode, err := vfs.resolvePath(originalPath)
	if err != nil {
		return &os.PathError{Op: "createvariant", Path: originalPath, Err: err}
	}

	if origInode.typ != "file" {
		return &os.PathError{Op: "createvariant", Path: originalPath, Err: errors.New("not a file")}
	}

	// Check if variant already exists
	var existingID int64
	err = vfs.db.db.QueryRow(
		"SELECT variant_inode_id FROM variants WHERE original_inode_id = ? AND variant_name = ?",
		origInode.id, variantName,
	).Scan(&existingID)

	if err == nil {
		return &os.PathError{Op: "createvariant", Path: originalPath, Err: os.ErrExist}
	}
	if err != sql.ErrNoRows {
		return err
	}

	tx, err := vfs.db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Create variant inode (hidden in normal directory listings)
	result, err := tx.Exec(
		"INSERT INTO inodes (fs_id, parent_id, name, type, mode, uid, gid, size, created_at, modified_at, accessed_at) VALUES (?, ?, ?, 'file', ?, ?, ?, ?, ?, ?, ?)",
		vfs.id, origInode.parentID, variantName, 0644, 0, 0, int64(len(content)), now, now, now,
	)
	if err != nil {
		return err
	}

	variantInodeID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Store variant content
	_, err = tx.Exec(
		"INSERT INTO storage (inode_id, content) VALUES (?, ?)",
		variantInodeID, content,
	)
	if err != nil {
		return err
	}

	// Create variant relationship
	_, err = tx.Exec(
		"INSERT INTO variants (original_inode_id, variant_inode_id, variant_name, created_at) VALUES (?, ?, ?, ?)",
		origInode.id, variantInodeID, variantName, now,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Emit variant event
	vfs.emitEvent(Event{
		Type:      EventVariant,
		Path:      originalPath + "/vo/" + variantName,
		Timestamp: now,
		FileType:  FileTypeFile,
		Size:      int64(len(content)),
	})

	return nil
}

// GetVariant opens a variant file for reading
func (vfs *VFS) GetVariant(originalPath string, variantName string) (*File, error) {
	originalPath = normalizePath(originalPath)

	// Get original file inode
	origInode, err := vfs.resolvePath(originalPath)
	if err != nil {
		return nil, &os.PathError{Op: "getvariant", Path: originalPath, Err: err}
	}

	// Get variant inode ID
	var variantInodeID int64
	err = vfs.db.db.QueryRow(
		"SELECT variant_inode_id FROM variants WHERE original_inode_id = ? AND variant_name = ?",
		origInode.id, variantName,
	).Scan(&variantInodeID)

	if err == sql.ErrNoRows {
		return nil, &os.PathError{Op: "getvariant", Path: originalPath, Err: os.ErrNotExist}
	}
	if err != nil {
		return nil, err
	}

	// Load variant inode
	variantInode, err := vfs.getInodeByID(variantInodeID)
	if err != nil {
		return nil, err
	}

	// Load content
	var content []byte
	err = vfs.db.db.QueryRow(
		"SELECT content FROM storage WHERE inode_id = ?",
		variantInodeID,
	).Scan(&content)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &File{
		vfs:      vfs,
		inode:    variantInode,
		path:     originalPath + "/vo/" + variantName,
		content:  content,
		pos:      0,
		flags:    os.O_RDONLY,
		writable: false,
		dirty:    false,
		closed:   false,
	}, nil
}

// ListVariants lists all variants of a file
func (vfs *VFS) ListVariants(originalPath string) ([]string, error) {
	originalPath = normalizePath(originalPath)

	// Get original file inode
	origInode, err := vfs.resolvePath(originalPath)
	if err != nil {
		return nil, &os.PathError{Op: "listvariants", Path: originalPath, Err: err}
	}

	rows, err := vfs.db.db.Query(
		"SELECT variant_name FROM variants WHERE original_inode_id = ? ORDER BY variant_name",
		origInode.id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		variants = append(variants, name)
	}

	return variants, rows.Err()
}

// RemoveVariant removes a variant of a file
func (vfs *VFS) RemoveVariant(originalPath string, variantName string) error {
	originalPath = normalizePath(originalPath)

	// Get original file inode
	origInode, err := vfs.resolvePath(originalPath)
	if err != nil {
		return &os.PathError{Op: "removevariant", Path: originalPath, Err: err}
	}

	// Delete variant (cascade will handle inode and storage)
	result, err := vfs.db.db.Exec(
		"DELETE FROM variants WHERE original_inode_id = ? AND variant_name = ?",
		origInode.id, variantName,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &os.PathError{Op: "removevariant", Path: originalPath, Err: os.ErrNotExist}
	}

	// Emit variant event
	vfs.emitEvent(Event{
		Type:      EventVariant,
		Path:      originalPath + "/vo/" + variantName,
		Timestamp: currentTimestamp(),
		FileType:  FileTypeFile,
	})

	return nil
}
