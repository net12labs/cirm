// filepath: /home/lxk/Desktop/cirm/bins/vfsql/types.go
package vfsql

import (
	"io"
	"io/fs"
	"os"
	"sync"
	"time"
)

// FileType represents the type of file system object
type FileType int

const (
	FileTypeAny FileType = iota
	FileTypeFile
	FileTypeDir
)

// SortField represents fields to sort by
type SortField int

const (
	SortByName SortField = iota
	SortBySize
	SortByCreated
	SortByModified
	SortByAccessed
	SortByPath
)

// SortOrder represents sort direction
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// EventType represents types of filesystem events
type EventType int

const (
	EventCreate EventType = iota
	EventModify
	EventDelete
	EventRename
	EventChmod
	EventChown
	EventMetadata
	EventVariant
)

// Event represents a filesystem event
type Event struct {
	Type      EventType
	Path      string
	OldPath   string
	Timestamp int64
	FileType  FileType
	Size      int64
	Tags      []string
}

// EventFilter specifies which events to receive
type EventFilter struct {
	Paths       []string
	Recursive   bool
	EventTypes  []EventType
	FileTypes   []FileType
	NamePattern string
	Tags        []string
	TagMatchAll bool
	BufferSize  int
}

// EventSubscription represents an active event subscription
type EventSubscription struct {
	ID     string
	Events chan Event
	Errors chan error
	filter *EventFilter
	cancel chan struct{}
}

// SearchQuery specifies search criteria
type SearchQuery struct {
	NamePattern    string
	NameRegex      string
	BasePath       string
	Recursive      bool
	MaxDepth       int
	Type           FileType
	Tags           []string
	TagMatchAll    bool
	Description    string
	MinSize        int64
	MaxSize        int64
	CreatedAfter   int64
	CreatedBefore  int64
	ModifiedAfter  int64
	ModifiedBefore int64
	AccessedAfter  int64
	AccessedBefore int64
	Mode           os.FileMode
	UID            int
	GID            int
	Limit          int
	Offset         int
	SortBy         SortField
	SortOrder      SortOrder
}

// SearchResults contains search results
type SearchResults struct {
	Paths      []string
	TotalCount int
	HasMore    bool
}

// FindOptions specifies options for simple find operations
type FindOptions struct {
	BasePath  string
	Recursive bool
	Type      FileType
	Limit     int
}

// VFS represents a virtual file system
type VFS struct {
	db          *DB
	volume      *Volume
	id          int64
	name        string
	rootInodeID int64
	cwd         string

	// Event subscription management
	subscribers map[string]*EventSubscription
	subMutex    sync.RWMutex
}

// File represents an open file
type File struct {
	vfs      *VFS
	inode    *inode
	path     string
	content  []byte
	pos      int64
	flags    int
	writable bool
	dirty    bool
	closed   bool
}

// inode represents internal inode data
type inode struct {
	id          int64
	fsID        int64
	parentID    *int64
	name        string
	typ         string
	mode        os.FileMode
	uid         int
	gid         int
	size        int64
	description string
	tags        string
	createdAt   int64
	modifiedAt  int64
	accessedAt  int64
}

// fileInfo implements os.FileInfo
type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     *inode
}

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *fileInfo) ModTime() time.Time { return fi.modTime }
func (fi *fileInfo) IsDir() bool        { return fi.isDir }
func (fi *fileInfo) Sys() interface{}   { return fi.sys }

// dirEntry implements fs.DirEntry
type dirEntry struct {
	info *fileInfo
}

func (de *dirEntry) Name() string               { return de.info.Name() }
func (de *dirEntry) IsDir() bool                { return de.info.IsDir() }
func (de *dirEntry) Type() fs.FileMode          { return de.info.Mode().Type() }
func (de *dirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

// Ensure interfaces are implemented
var (
	_ fs.FileInfo = (*fileInfo)(nil)
	_ fs.DirEntry = (*dirEntry)(nil)
	_ io.Reader   = (*File)(nil)
	_ io.Writer   = (*File)(nil)
	_ io.Closer   = (*File)(nil)
	_ io.Seeker   = (*File)(nil)
)
