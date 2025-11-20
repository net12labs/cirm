package vfsql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the SQLite database containing all volumes and filesystems
type DB struct {
	db   *sql.DB
	path string
}

// Open opens an existing database or creates a new one if it doesn't exist
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	// Initialize schema if needed
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db: db, path: path}, nil
}

// Create creates a new database file
func Create(path string) (*DB, error) {
	// Remove existing file if present
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}

	return Open(path)
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// initSchema creates all necessary tables if they don't exist
func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS volumes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS filesystems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		volume_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		root_inode_id INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		FOREIGN KEY (volume_id) REFERENCES volumes(id) ON DELETE CASCADE,
		FOREIGN KEY (root_inode_id) REFERENCES inodes(id) ON DELETE SET NULL,
		UNIQUE(volume_id, name)
	);

	CREATE INDEX IF NOT EXISTS idx_filesystems_volume ON filesystems(volume_id);

	CREATE TABLE IF NOT EXISTS inodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fs_id INTEGER NOT NULL,
		parent_id INTEGER,
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('file', 'dir')),
		mode INTEGER NOT NULL,
		uid INTEGER NOT NULL,
		gid INTEGER NOT NULL,
		size INTEGER NOT NULL DEFAULT 0,
		description TEXT,
		tags TEXT,
		created_at INTEGER NOT NULL,
		modified_at INTEGER NOT NULL,
		accessed_at INTEGER NOT NULL,
		FOREIGN KEY (fs_id) REFERENCES filesystems(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES inodes(id) ON DELETE CASCADE,
		UNIQUE(fs_id, parent_id, name)
	);

	CREATE INDEX IF NOT EXISTS idx_inodes_fs_parent ON inodes(fs_id, parent_id);
	CREATE INDEX IF NOT EXISTS idx_inodes_path ON inodes(fs_id, parent_id, name);
	CREATE INDEX IF NOT EXISTS idx_inodes_tags ON inodes(tags);

	CREATE TABLE IF NOT EXISTS storage (
		inode_id INTEGER PRIMARY KEY,
		content BLOB NOT NULL,
		FOREIGN KEY (inode_id) REFERENCES inodes(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS variants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_inode_id INTEGER NOT NULL,
		variant_inode_id INTEGER NOT NULL,
		variant_name TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		FOREIGN KEY (original_inode_id) REFERENCES inodes(id) ON DELETE CASCADE,
		FOREIGN KEY (variant_inode_id) REFERENCES inodes(id) ON DELETE CASCADE,
		UNIQUE(original_inode_id, variant_name)
	);

	CREATE INDEX IF NOT EXISTS idx_variants_original ON variants(original_inode_id);
	CREATE INDEX IF NOT EXISTS idx_variants_variant ON variants(variant_inode_id);
	`

	_, err := db.Exec(schema)
	return err
}

// CreateVolume creates a new volume with a default filesystem
func (d *DB) CreateVolume(name string) (*Volume, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Insert volume
	result, err := tx.Exec(
		"INSERT INTO volumes (name, created_at, updated_at) VALUES (?, ?, ?)",
		name, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume: %w", err)
	}

	volumeID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create default filesystem
	fsResult, err := tx.Exec(
		"INSERT INTO filesystems (volume_id, name, created_at, updated_at) VALUES (?, 'default', ?, ?)",
		volumeID, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create default filesystem: %w", err)
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

	return &Volume{
		db:   d,
		id:   volumeID,
		name: name,
	}, nil
}

// GetVolume retrieves an existing volume by name
func (d *DB) GetVolume(name string) (*Volume, error) {
	var id int64
	err := d.db.QueryRow("SELECT id FROM volumes WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, os.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	return &Volume{
		db:   d,
		id:   id,
		name: name,
	}, nil
}

// DeleteVolume deletes a volume and all its data
func (d *DB) DeleteVolume(name string) error {
	result, err := d.db.Exec("DELETE FROM volumes WHERE name = ?", name)
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

// ListVolumes returns all volume names
func (d *DB) ListVolumes() ([]string, error) {
	rows, err := d.db.Query("SELECT name FROM volumes ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		volumes = append(volumes, name)
	}

	return volumes, rows.Err()
}
