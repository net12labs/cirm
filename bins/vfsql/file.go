// filepath: /home/lxk/Desktop/cirm/bins/vfsql/file.go
package vfsql

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

// Read reads up to len(p) bytes from the file
func (f *File) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if f.inode.typ != "file" {
		return 0, &os.PathError{Op: "read", Path: f.path, Err: errors.New("is a directory")}
	}

	if f.pos >= int64(len(f.content)) {
		return 0, io.EOF
	}

	n := copy(p, f.content[f.pos:])
	f.pos += int64(n)

	return n, nil
}

// Write writes len(p) bytes to the file
func (f *File) Write(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if !f.writable {
		return 0, &os.PathError{Op: "write", Path: f.path, Err: os.ErrPermission}
	}

	if f.inode.typ != "file" {
		return 0, &os.PathError{Op: "write", Path: f.path, Err: errors.New("is a directory")}
	}

	// Expand buffer if necessary
	minLen := int(f.pos) + len(p)
	if minLen > len(f.content) {
		newContent := make([]byte, minLen)
		copy(newContent, f.content)
		f.content = newContent
	}

	n := copy(f.content[f.pos:], p)
	f.pos += int64(n)
	f.dirty = true

	// Update size
	if f.pos > f.inode.size {
		f.inode.size = f.pos
	}

	fmt.Printf("[File.Write] inode=%d, wrote %d bytes, new size=%d, dirty=%v\n",
		f.inode.id, n, f.inode.size, f.dirty)

	return n, nil
}

// Seek sets the offset for the next Read or Write
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = f.pos + offset
	case io.SeekEnd:
		newPos = int64(len(f.content)) + offset
	default:
		return 0, &os.PathError{Op: "seek", Path: f.path, Err: errors.New("invalid whence")}
	}

	if newPos < 0 {
		return 0, &os.PathError{Op: "seek", Path: f.path, Err: errors.New("negative position")}
	}

	f.pos = newPos
	return newPos, nil
}

// Close closes the file, writing any changes to the database
func (f *File) Close() error {
	if f.closed {
		return fs.ErrClosed
	}

	fmt.Printf("[File.Close] inode=%d, content len=%d, dirty=%v\n",
		f.inode.id, len(f.content), f.dirty)

	// Write changes if dirty BEFORE marking as closed
	if f.dirty && f.writable {
		if err := f.Sync(); err != nil {
			fmt.Printf("[File.Close] Sync failed: %v\n", err)
			f.closed = true
			return err
		}
	}

	f.closed = true
	return nil
}

// Sync writes any buffered changes to the database
func (f *File) Sync() error {
	if f.closed {
		return fs.ErrClosed
	}

	if !f.dirty || !f.writable {
		return nil
	}

	fmt.Printf("[File.Sync] inode=%d, size=%d, content_len=%d\n",
		f.inode.id, f.inode.size, len(f.content))

	tx, err := f.vfs.db.db.Begin()
	if err != nil {
		fmt.Printf("[File.Sync] Begin transaction failed: %v\n", err)
		return err
	}
	defer tx.Rollback()

	now := currentTimestamp()

	// Update inode size and modified time
	result, err := tx.Exec(
		"UPDATE inodes SET size = ?, modified_at = ? WHERE id = ?",
		f.inode.size, now, f.inode.id,
	)
	if err != nil {
		fmt.Printf("[File.Sync] UPDATE inodes failed: %v\n", err)
		return err
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("[File.Sync] UPDATE inodes affected %d rows\n", rows)

	// Update storage content - store exactly the content we have
	contentToStore := f.content[:f.inode.size]
	result, err = tx.Exec(
		"INSERT OR REPLACE INTO storage (inode_id, content) VALUES (?, ?)",
		f.inode.id, contentToStore,
	)
	if err != nil {
		fmt.Printf("[File.Sync] INSERT storage failed: %v\n", err)
		return err
	}
	rows, _ = result.RowsAffected()
	fmt.Printf("[File.Sync] INSERT storage affected %d rows, stored %d bytes\n", rows, len(contentToStore))

	if err := tx.Commit(); err != nil {
		fmt.Printf("[File.Sync] Commit failed: %v\n", err)
		return err
	}

	fmt.Printf("[File.Sync] SUCCESS - committed %d bytes to storage\n", len(contentToStore))

	f.dirty = false
	f.inode.modifiedAt = now

	// Emit modify event
	f.vfs.emitEvent(Event{
		Type:      EventModify,
		Path:      f.path,
		Timestamp: now,
		FileType:  FileTypeFile,
		Size:      f.inode.size,
	})

	return nil
}

// Name returns the name of the file
func (f *File) Name() string {
	return f.path
}

// Stat returns file info
func (f *File) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}

	return &fileInfo{
		name:    f.inode.name,
		size:    f.inode.size,
		mode:    f.inode.mode,
		modTime: time.Unix(f.inode.modifiedAt, 0),
		isDir:   f.inode.typ == "dir",
		sys:     f.inode,
	}, nil
}

// Readdir reads directory entries
func (f *File) Readdir(n int) ([]fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}

	if f.inode.typ != "dir" {
		return nil, &os.PathError{Op: "readdir", Path: f.path, Err: errors.New("not a directory")}
	}

	entries, err := f.vfs.ReadDir(f.path)
	if err != nil {
		return nil, err
	}

	infos := make([]fs.FileInfo, len(entries))
	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos[i] = info
	}

	if n <= 0 {
		return infos, nil
	}

	if n > len(infos) {
		n = len(infos)
	}

	return infos[:n], nil
}

// Readdirnames reads directory entry names
func (f *File) Readdirnames(n int) ([]string, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}

	if f.inode.typ != "dir" {
		return nil, &os.PathError{Op: "readdirnames", Path: f.path, Err: errors.New("not a directory")}
	}

	entries, err := f.vfs.ReadDir(f.path)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}

	if n <= 0 {
		return names, nil
	}

	if n > len(names) {
		n = len(names)
	}

	return names[:n], nil
}

// Truncate changes the size of the file
func (f *File) Truncate(size int64) error {
	if f.closed {
		return fs.ErrClosed
	}

	if !f.writable {
		return &os.PathError{Op: "truncate", Path: f.path, Err: os.ErrPermission}
	}

	if f.inode.typ != "file" {
		return &os.PathError{Op: "truncate", Path: f.path, Err: errors.New("is a directory")}
	}

	if size < 0 {
		return &os.PathError{Op: "truncate", Path: f.path, Err: errors.New("negative size")}
	}

	if size > int64(len(f.content)) {
		newContent := make([]byte, size)
		copy(newContent, f.content)
		f.content = newContent
	} else {
		f.content = f.content[:size]
	}

	f.inode.size = size
	f.pos = 0
	f.dirty = true

	return nil
}
