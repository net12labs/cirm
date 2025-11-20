// filepath: /home/lxk/Desktop/cirm/bins/vfsql/volume.go
package vfsql

import (
	"database/sql"
	"fmt"
	"os"
)

// Volume represents a logical partition within the database
type Volume struct {
	db   *DB
	id   int64
	name string
}

// CreateVFS creates a new virtual filesystem within the volume
func (v *Volume) CreateVFS(name string) (*VFS, error) {
	tx, err := v.db.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Insert filesystem
	fsResult, err := tx.Exec(
		"INSERT INTO filesystems (volume_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		v.id, name, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem: %w", err)
	}

	fsID, err := fsResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create root directory inode
	rootResult, err := tx.Exec(
		"INSERT INTO inodes (fs_id, parent_id, name, type, mode, uid, gid, created_at, modified_at, accessed_at) VALUES (?, NULL, '/', 'dir', ?, ?, ?, ?, ?, ?)",
		fsID, 0755, 0, 0, now, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create root inode: %w", err)
	}

	rootInodeID, err := rootResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update filesystem with root_inode_id
	_, err = tx.Exec(
		"UPDATE filesystems SET root_inode_id = ? WHERE id = ?",
		rootInodeID, fsID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update filesystem root: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &VFS{
		db:          v.db,
		volume:      v,
		id:          fsID,
		name:        name,
		rootInodeID: rootInodeID,
		cwd:         "/",
		subscribers: make(map[string]*EventSubscription),
	}, nil
}

// GetVFS retrieves an existing VFS by name
func (v *Volume) GetVFS(name string) (*VFS, error) {
	var id, rootInodeID int64
	err := v.db.db.QueryRow(
		"SELECT id, root_inode_id FROM filesystems WHERE volume_id = ? AND name = ?",
		v.id, name,
	).Scan(&id, &rootInodeID)

	if err == sql.ErrNoRows {
		return nil, os.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	return &VFS{
		db:          v.db,
		volume:      v,
		id:          id,
		name:        name,
		rootInodeID: rootInodeID,
		cwd:         "/",
		subscribers: make(map[string]*EventSubscription),
	}, nil
}

// DeleteVFS deletes a virtual filesystem
func (v *Volume) DeleteVFS(name string) error {
	result, err := v.db.db.Exec(
		"DELETE FROM filesystems WHERE volume_id = ? AND name = ?",
		v.id, name,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return os.ErrNotExist
	}

	return nil
}

// ListVFS returns all filesystem names in the volume
func (v *Volume) ListVFS() ([]string, error) {
	rows, err := v.db.db.Query(
		"SELECT name FROM filesystems WHERE volume_id = ? ORDER BY name",
		v.id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filesystems []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		filesystems = append(filesystems, name)
	}

	return filesystems, rows.Err()
}
