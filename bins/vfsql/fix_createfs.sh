#!/bin/bash
# Fix CreateFilesystem function in volume.go

cd /home/lxk/Desktop/cirm/bins/vfsql

# Find the line number where CreateFilesystem starts
START=$(grep -n "^func (v \*Volume) CreateFilesystem" volume.go | cut -d: -f1)

if [ -z "$START" ]; then
    echo "CreateFilesystem function not found"
    exit 1
fi

# Find the end of the first occurrence (find matching closing brace)
# This is a simple approach - find the next function or end of file
END=$(tail -n +$((START+1)) volume.go | grep -n "^func \|^//.*\n^func \|^$" | head -1 | cut -d: -f1)
END=$((START + END))

echo "Found CreateFilesystem from line $START to $END"
echo "Backing up to volume.go.bak2"
cp volume.go volume.go.bak2

# Create new version with proper transaction handling
cat > /tmp/createfs.go << 'EOFUNC'
// CreateFilesystem creates a new filesystem in this volume
func (v *Volume) CreateFilesystem(name string) error {
	tx, err := v.db.db.Begin()
	if err != nil {
		return err
	}
	
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	
	now := currentTimestamp()
	
	// Create filesystem
	fsResult, err := tx.Exec(
		"INSERT INTO filesystems (volume_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		v.id, name, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create filesystem: %w", err)
	}
	
	fsID, err := fsResult.LastInsertId()
	if err != nil {
		return err
	}
	
	// Create root directory inode
	rootResult, err := tx.Exec(
		"INSERT INTO inodes (fs_id, parent_id, name, type, mode, uid, gid, created_at, modified_at, accessed_at) VALUES (?, NULL, '/', 'dir', ?, ?, ?, ?, ?, ?)",
		fsID, 0755, 0, 0, now, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create root inode: %w", err)
	}
	
	rootInodeID, err := rootResult.LastInsertId()
	if err != nil {
		return err
	}
	
	// Update filesystem with root_inode_id
	_, err = tx.Exec(
		"UPDATE filesystems SET root_inode_id = ? WHERE id = ?",
		rootInodeID, fsID,
	)
	if err != nil {
		return fmt.Errorf("failed to update filesystem root: %w", err)
	}
	
	err = tx.Commit()
	if err != nil {
		return err
	}
	
	committed = true
	return nil
}
EOFUNC

# Reconstruct file: before function + new function + after function  
head -n $((START-1)) volume.go > volume.go.new
cat /tmp/createfs.go >> volume.go.new
tail -n +$END volume.go >> volume.go.new

mv volume.go.new volume.go
echo "✓ Fixed CreateFilesystem function"
