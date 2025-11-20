# Virtual File System (VFSQL) Specification

## Overview
A virtual file system backed by a **single SQLite database file** that provides os-compatible file access patterns in Go. Designed for small files without chunking requirements.

All volumes, virtual file systems, directories, and files are stored in one SQLite database.

## Architecture

### Hierarchy
```
SQLite Database → Volumes → VFS → Directories → Files
```

### Components

#### 0. SQLite Database
- Single .db file containing everything
- All volumes, VFS instances, directories, and files stored together
- Shared connection pool and transaction management

#### 1. Volume
- Logical partition within the database
- Multiple volumes can exist in the same database
- Provides namespace isolation between different volume instances
- Identified by unique name
- **Automatically creates a default filesystem with root directory on initialization**

#### 2. VFS (Virtual File System)
- Named file system within a volume
- Has its own root directory inode
- Multiple VFS instances can exist per volume
- Provides additional namespace isolation
- Root directory (/) is automatically created when filesystem is initialized

#### 3. Directories
- Hierarchical structure within a VFS
- Can contain files and subdirectories
- Standard path operations (/, separators)

#### 4. Files
- Small file storage (no chunking)
- Metadata (inode) separate from content (storage)

---

## Database Schema

### Table: volumes
```sql
CREATE TABLE volumes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
```

### Table: filesystems
```sql
CREATE TABLE filesystems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    volume_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    root_inode_id INTEGER, -- Reference to root directory inode
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (volume_id) REFERENCES volumes(id) ON DELETE CASCADE,
    FOREIGN KEY (root_inode_id) REFERENCES inodes(id) ON DELETE SET NULL,
    UNIQUE(volume_id, name)
);

CREATE INDEX idx_filesystems_volume ON filesystems(volume_id);
```

### Table: inodes
Stores file and directory metadata
```sql
CREATE TABLE inodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fs_id INTEGER NOT NULL,
    parent_id INTEGER, -- NULL for root
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('file', 'dir')),
    mode INTEGER NOT NULL, -- Unix permissions
    uid INTEGER NOT NULL,
    gid INTEGER NOT NULL,
    size INTEGER NOT NULL DEFAULT 0,
    description TEXT, -- Optional description
    tags TEXT, -- Comma-separated tags for search/categorization
    created_at INTEGER NOT NULL,
    modified_at INTEGER NOT NULL,
    accessed_at INTEGER NOT NULL,
    FOREIGN KEY (fs_id) REFERENCES filesystems(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES inodes(id) ON DELETE CASCADE,
    UNIQUE(fs_id, parent_id, name)
);

CREATE INDEX idx_inodes_fs_parent ON inodes(fs_id, parent_id);
CREATE INDEX idx_inodes_path ON inodes(fs_id, parent_id, name);
CREATE INDEX idx_inodes_tags ON inodes(tags); -- For tag-based searches
```

### Table: storage
Stores actual file content
```sql
CREATE TABLE storage (
    inode_id INTEGER PRIMARY KEY,
    content BLOB NOT NULL,
    FOREIGN KEY (inode_id) REFERENCES inodes(id) ON DELETE CASCADE
);
```

### Table: variants
Stores variant relationships between files
```sql
CREATE TABLE variants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_inode_id INTEGER NOT NULL,
    variant_inode_id INTEGER NOT NULL,
    variant_name TEXT NOT NULL, -- e.g., "200x300px.gif", "thumbnail.jpg"
    created_at INTEGER NOT NULL,
    FOREIGN KEY (original_inode_id) REFERENCES inodes(id) ON DELETE CASCADE,
    FOREIGN KEY (variant_inode_id) REFERENCES inodes(id) ON DELETE CASCADE,
    UNIQUE(original_inode_id, variant_name)
);

CREATE INDEX idx_variants_original ON variants(original_inode_id);
CREATE INDEX idx_variants_variant ON variants(variant_inode_id);
```

---

## Go API Design

### Core Interfaces

#### VFS Interface
```go
type VFS interface {
    // File operations
    Open(path string) (File, error)
    OpenFile(path string, flag int, perm os.FileMode) (File, error)
    Create(path string) (File, error)
    Remove(path string) error
    Rename(oldpath, newpath string) error
    
    // Directory operations
    Mkdir(path string, perm os.FileMode) error
    MkdirAll(path string, perm os.FileMode) error
    RemoveAll(path string) error
    ReadDir(path string) ([]DirEntry, error)
    
    // Stat operations
    Stat(path string) (FileInfo, error)
    Lstat(path string) (FileInfo, error)
    
    // Permissions
    Chmod(path string, mode os.FileMode) error
    Chown(path string, uid, gid int) error
    
    // Working directory
    Getwd() (string, error)
    Chdir(path string) error
    
    // Extended metadata (VFSQL-specific)
    SetDescription(path string, description string) error
    GetDescription(path string) (string, error)
    SetTags(path string, tags []string) error
    GetTags(path string) ([]string, error)
    AddTag(path string, tag string) error
    RemoveTag(path string, tag string) error
    
    // Search operations
    Search(query *SearchQuery) (*SearchResults, error)
    FindByName(pattern string, opts *FindOptions) ([]string, error)
    FindByTag(tags []string, matchAll bool) ([]string, error)
    
    // Variant management
    CreateVariant(originalPath string, variantName string, content []byte) error
    GetVariant(originalPath string, variantName string) (File, error)
    ListVariants(originalPath string) ([]string, error)
    RemoveVariant(originalPath string, variantName string) error
    
    // Event subscription
    Subscribe(filter *EventFilter) (*EventSubscription, error)
    Unsubscribe(sub *EventSubscription) error
}
```

#### SearchQuery Structure
```go
type SearchQuery struct {
    // Name pattern matching
    NamePattern string // Glob pattern (*, ?, [abc]) or empty
    NameRegex   string // Regex pattern or empty
    
    // Path constraints
    BasePath  string // Search only within this path
    Recursive bool   // Search subdirectories
    MaxDepth  int    // Maximum depth (0 = unlimited)
    
    // Type filter
    Type FileType // Filter by type: FileTypeAny, FileTypeFile, FileTypeDir
    
    // Metadata filters
    Tags        []string // Tag filters
    TagMatchAll bool     // true = AND, false = OR
    Description string   // Substring match in description
    
    // Size filters
    MinSize int64 // Minimum file size in bytes (0 = no limit)
    MaxSize int64 // Maximum file size in bytes (0 = no limit)
    
    // Time filters (Unix timestamps)
    CreatedAfter   int64
    CreatedBefore  int64
    ModifiedAfter  int64
    ModifiedBefore int64
    AccessedAfter  int64
    AccessedBefore int64
    
    // Permission filters
    Mode os.FileMode // Filter by permission bits (0 = ignore)
    UID  int         // Filter by user ID (-1 = ignore)
    GID  int         // Filter by group ID (-1 = ignore)
    
    // Result options
    Limit      int  // Maximum results (0 = unlimited)
    Offset     int  // Skip first N results
    SortBy     SortField
    SortOrder  SortOrder // Ascending or Descending
}

type FileType int
const (
    FileTypeAny FileType = iota
    FileTypeFile
    FileTypeDir
)

type SortField int
const (
    SortByName SortField = iota
    SortBySize
    SortByCreated
    SortByModified
    SortByAccessed
    SortByPath
)

type SortOrder int
const (
    Ascending SortOrder = iota
    Descending
)

type SearchResults struct {
    Paths      []string
    TotalCount int // Total matches (before limit/offset)
    HasMore    bool
}
```

#### FindOptions Structure
```go
type FindOptions struct {
    BasePath  string
    Recursive bool
    Type      FileType
    Limit     int
}
```

#### Event System Structures
```go
type EventType int
const (
    EventCreate EventType = iota
    EventModify
    EventDelete
    EventRename
    EventChmod
    EventChown
    EventMetadata // Description or tags changed
    EventVariant  // Variant created/deleted
)

type Event struct {
    Type      EventType
    Path      string
    OldPath   string    // For rename operations
    Timestamp int64     // Unix timestamp
    FileType  FileType  // File or directory
    Size      int64     // File size (for create/modify)
    Tags      []string  // Current tags (for metadata events)
}

type EventFilter struct {
    // Path filtering
    Paths     []string // Watch specific paths
    Recursive bool     // Watch subdirectories
    
    // Event type filtering
    EventTypes []EventType // Empty = all events
    
    // File type filtering
    FileTypes []FileType // Empty = all types
    
    // Pattern matching
    NamePattern string // Glob pattern for file names
    
    // Tag filtering
    Tags        []string // Only notify for files with these tags
    TagMatchAll bool     // true = AND, false = OR
    
    // Buffering
    BufferSize int // Event buffer size (default: 100)
}

type EventSubscription struct {
    ID     string
    Events <-chan Event
    Errors <-chan error
    filter *EventFilter
}
```

#### File Interface
```go
type File interface {
    io.Reader
    io.Writer
    io.Closer
    io.Seeker
    
    Name() string
    Stat() (FileInfo, error)
    Readdir(n int) ([]FileInfo, error)
    Readdirnames(n int) ([]string, error)
    Truncate(size int64) error
    Sync() error
}
```

#### FileInfo Interface (compatible with os.FileInfo)
```go
type FileInfo interface {
    Name() string
    Size() int64
    Mode() os.FileMode
    ModTime() time.Time
    IsDir() bool
    Sys() interface{}
}
```

#### DirEntry Interface (compatible with os.DirEntry)
```go
type DirEntry interface {
    Name() string
    IsDir() bool
    Type() os.FileMode
    Info() (FileInfo, error)
}
```

### Database Manager

```go
type DB struct {
    db *sql.DB
    path string
}

// Database lifecycle
func Open(path string) (*DB, error)
func Create(path string) (*DB, error)
func (d *DB) Close() error

// Volume operations
func (d *DB) CreateVolume(name string) (*Volume, error)
func (d *DB) GetVolume(name string) (*Volume, error)
func (d *DB) DeleteVolume(name string) error
func (d *DB) ListVolumes() ([]string, error)

// Volume represents a logical partition
type Volume struct {
    db *DB
    id int64
    name string
}

// VFS operations on a volume
func (v *Volume) CreateVFS(name string) (*VFS, error)
func (v *Volume) GetVFS(name string) (*VFS, error)
func (v *Volume) DeleteVFS(name string) error
func (v *Volume) ListVFS() ([]string, error)
```

---

## Implementation Details

### Path Resolution
- Normalize paths (handle `.`, `..`, multiple slashes)
- Resolve absolute and relative paths
- Validate path components
- Track current working directory per VFS instance
- **Special path handling**: `/vo/` subdirectory pattern for file variants
  - Pattern: `<filepath>/vo/<variant-name>`
  - Example: `/images/photo.jpg/vo/200x300px.gif`
  - The `/vo/` is a virtual directory that doesn't exist as a real inode
  - Accessing paths with `/vo/` triggers variant lookup logic

### File Operations
- **Open modes**: Read, Write, Append, Create, Exclusive, Truncate
- **File descriptor**: In-memory buffer for read/write operations
- **Sync**: Write buffer to SQLite on Sync() or Close()
- **Concurrent access**: Use transactions for consistency

### Initialization Sequence
When a new volume is created:
1. Insert volume record into `volumes` table
2. Create default filesystem named "default" in `filesystems` table
3. Create root directory inode (parent_id = NULL, name = "/", type = "dir", mode = 0755)
4. Update filesystem record with `root_inode_id`

When a new filesystem is created within a volume:
1. Insert filesystem record into `filesystems` table
2. Create root directory inode for the filesystem
3. Update filesystem record with `root_inode_id`

### Directory Operations
- Recursive directory creation (MkdirAll)
- Recursive deletion (RemoveAll)
- List directory contents with filtering
- Prevent deletion of non-empty directories (for Remove)
- Root directory (/) cannot be deleted

### Permissions
- Store Unix-style permissions (rwxrwxrwx)
- Default: 0644 for files, 0755 for directories
- UID/GID support (may be stub values on Windows)

### Timestamps
- created_at: Creation time (immutable)
- modified_at: Last modification time
- accessed_at: Last access time
- Store as Unix timestamps (seconds since epoch)

### Error Handling
- Return standard Go errors compatible with os package
- Use os.PathError for path-related errors
- Support errors.Is() checks for common errors:
  - os.ErrNotExist
  - os.ErrExist
  - os.ErrPermission
  - os.ErrInvalid

### Transaction Strategy
- Read operations: Use shared reads
- Write operations: Wrap in transactions
- File writes: Defer actual storage write until Sync/Close
- Batch operations: Single transaction when possible

### Search Implementation
- **Query Building**: Dynamically construct SQL WHERE clauses based on SearchQuery
- **Path reconstruction**: Use recursive CTE or iterative parent_id lookup to build full paths
- **Tag matching**: 
  - OR logic: Use `tags LIKE '%tag1%' OR tags LIKE '%tag2%'`
  - AND logic: Ensure all tags present (multiple LIKE conditions)
- **Pattern matching**:
  - Glob patterns: Convert to SQL LIKE (e.g., `*.txt` → `%.txt`)
  - Regex: Use SQLite REGEXP function (may require extension)
- **Optimization**:
  - Use indexes on (fs_id, parent_id, name) for path traversal
  - Use tag index for tag-based searches
  - Use covering indexes where possible
  - COUNT query for TotalCount separate from data query
- **Pagination**: Use LIMIT and OFFSET for result paging
- **Sorting**: ORDER BY clause based on SortBy and SortOrder

### Event System Implementation
- **Event Generation**: All write operations (Create, Modify, Delete, etc.) emit events
- **Subscription Management**:
  - Keep registry of active subscriptions per VFS
  - Each subscription has unique ID
  - Track filter criteria for each subscription
- **Event Dispatch**:
  - After successful operation, check all active subscriptions
  - Match event against each subscription's filter
  - Send event to matching subscriptions via buffered channels
  - Non-blocking send: drop events if channel is full (with error notification)
- **Filter Matching**:
  - Path matching: Check if event path is within subscription paths
  - Recursive: Check if event is in subdirectories
  - Pattern matching: Apply glob pattern to file name
  - Tag matching: Check if file's tags match subscription tags
  - Event type filtering: Only dispatch matching event types
- **Thread Safety**: Use mutex to protect subscription registry
- **Cleanup**: Automatically clean up closed subscriptions
- **Performance**: 
  - Avoid blocking write operations
  - Use goroutines for event dispatch
  - Buffer events to handle temporary slowness in consumers

---

## Usage Examples

### Creating a Database with Volumes and VFS
```go
// Open or create the database
db, err := vfsql.Open("/tmp/mydata.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Create a volume
vol, err := db.CreateVolume("myvolume")
if err != nil {
    log.Fatal(err)
}

// Create a VFS within the volume
vfs, err := vol.CreateVFS("default")
if err != nil {
    log.Fatal(err)
}
```

### File Operations
```go
// Create a file
f, err := vfs.Create("/hello.txt")
if err != nil {
    log.Fatal(err)
}
f.Write([]byte("Hello, World!"))
f.Close()

// Read a file
f, err = vfs.Open("/hello.txt")
if err != nil {
    log.Fatal(err)
}
content, _ := io.ReadAll(f)
f.Close()
fmt.Println(string(content))
```

### Directory Operations
```go
// Create directories
err := vfs.MkdirAll("/path/to/dir", 0755)

// List directory
entries, err := vfs.ReadDir("/path/to")
for _, entry := range entries {
    fmt.Println(entry.Name(), entry.IsDir())
}

// Remove directory tree
err = vfs.RemoveAll("/path")
```

### Extended Metadata Operations
```go
// Set description on a file
err := vfs.SetDescription("/hello.txt", "My first file")

// Get description
desc, err := vfs.GetDescription("/hello.txt")

// Set tags
err = vfs.SetTags("/hello.txt", []string{"important", "draft", "2024"})

// Add a single tag
err = vfs.AddTag("/hello.txt", "reviewed")

// Get all tags
tags, err := vfs.GetTags("/hello.txt")

// Remove a tag
err = vfs.RemoveTag("/hello.txt", "draft")
```

### Search Operations
```go
// Simple name-based search
paths, err := vfs.FindByName("*.jpg", &vfsql.FindOptions{
    BasePath:  "/images",
    Recursive: true,
    Type:      vfsql.FileTypeFile,
})

// Search by tags
paths, err = vfs.FindByTag([]string{"important", "reviewed"}, true) // true = match all

// Advanced search with multiple criteria
results, err := vfs.Search(&vfsql.SearchQuery{
    // Name filtering
    NamePattern: "report*.pdf",
    
    // Search in specific directory
    BasePath:  "/documents",
    Recursive: true,
    MaxDepth:  3,
    
    // Type and metadata filters
    Type:        vfsql.FileTypeFile,
    Tags:        []string{"2024", "financial"},
    TagMatchAll: true,
    Description: "quarterly",
    
    // Size constraints
    MinSize: 1024,        // At least 1KB
    MaxSize: 10485760,    // At most 10MB
    
    // Time range (last 30 days)
    ModifiedAfter: time.Now().AddDate(0, 0, -30).Unix(),
    
    // Pagination and sorting
    Limit:     50,
    Offset:    0,
    SortBy:    vfsql.SortByModified,
    SortOrder: vfsql.Descending,
})

fmt.Printf("Found %d matches (showing %d)\n", results.TotalCount, len(results.Paths))
for _, path := range results.Paths {
    fmt.Println(path)
}
if results.HasMore {
    fmt.Println("More results available...")
}

// Search for large files modified recently
results, err = vfs.Search(&vfsql.SearchQuery{
    Type:          vfsql.FileTypeFile,
    MinSize:       5242880, // 5MB
    ModifiedAfter: time.Now().AddDate(0, 0, -7).Unix(),
    SortBy:        vfsql.SortBySize,
    SortOrder:     vfsql.Descending,
    Limit:         10,
})

// Find all directories with specific tag
results, err = vfs.Search(&vfsql.SearchQuery{
    Type: vfsql.FileTypeDir,
    Tags: []string{"project"},
    SortBy: vfsql.SortByName,
})

// Full-text search in descriptions
results, err = vfs.Search(&vfsql.SearchQuery{
    Description: "meeting notes",
    Recursive:   true,
    BasePath:    "/",
})
```

### File Variant Operations
```go
// Create an original image
f, err := vfs.Create("/images/photo.jpg")
f.Write(originalImageData)
f.Close()

// Create a variant (thumbnail)
thumbnailData := resizeImage(originalImageData, 200, 300)
err = vfs.CreateVariant("/images/photo.jpg", "200x300px.gif", thumbnailData)

// Access variant using virtual path
f, err = vfs.Open("/images/photo.jpg/vo/200x300px.gif")
content, _ := io.ReadAll(f)
f.Close()

// Or use the API directly
f, err = vfs.GetVariant("/images/photo.jpg", "200x300px.gif")

// List all variants of a file
variants, err := vfs.ListVariants("/images/photo.jpg")
// Returns: ["200x300px.gif", "thumbnail.jpg", ...]

// Remove a variant
err = vfs.RemoveVariant("/images/photo.jpg", "200x300px.gif")

// Reading directory shows original file, not variants
entries, err := vfs.ReadDir("/images")
// Returns: ["photo.jpg"] - variants are hidden from normal directory listings
```

### Event Subscription Operations
```go
// Subscribe to all events in a directory
sub, err := vfs.Subscribe(&vfsql.EventFilter{
    Paths:     []string{"/documents"},
    Recursive: true,
})
if err != nil {
    log.Fatal(err)
}
defer vfs.Unsubscribe(sub)

// Listen for events
go func() {
    for {
        select {
        case event := <-sub.Events:
            switch event.Type {
            case vfsql.EventCreate:
                fmt.Printf("File created: %s (size: %d)\n", event.Path, event.Size)
            case vfsql.EventModify:
                fmt.Printf("File modified: %s\n", event.Path)
            case vfsql.EventDelete:
                fmt.Printf("File deleted: %s\n", event.Path)
            case vfsql.EventRename:
                fmt.Printf("File renamed: %s -> %s\n", event.OldPath, event.Path)
            case vfsql.EventMetadata:
                fmt.Printf("Metadata changed: %s (tags: %v)\n", event.Path, event.Tags)
            }
        case err := <-sub.Errors:
            log.Printf("Event error: %v\n", err)
        }
    }
}()

// Subscribe to specific event types only
sub, err = vfs.Subscribe(&vfsql.EventFilter{
    Paths:      []string{"/"},
    Recursive:  true,
    EventTypes: []vfsql.EventType{vfsql.EventCreate, vfsql.EventDelete},
    FileTypes:  []vfsql.FileType{vfsql.FileTypeFile},
})

// Watch files with specific pattern
sub, err = vfs.Subscribe(&vfsql.EventFilter{
    Paths:       []string{"/images"},
    Recursive:   true,
    NamePattern: "*.jpg",
    EventTypes:  []vfsql.EventType{vfsql.EventCreate, vfsql.EventModify},
})

// Watch files with specific tags
sub, err = vfs.Subscribe(&vfsql.EventFilter{
    Paths:       []string{"/"},
    Recursive:   true,
    Tags:        []string{"important", "monitored"},
    TagMatchAll: true,
    BufferSize:  500, // Larger buffer for high-volume scenarios
})

// Multiple subscriptions can be active simultaneously
sub1, _ := vfs.Subscribe(&vfsql.EventFilter{Paths: []string{"/logs"}})
sub2, _ := vfs.Subscribe(&vfsql.EventFilter{Paths: []string{"/data"}})

// Process events from multiple subscriptions
for {
    select {
    case event := <-sub1.Events:
        handleLogEvent(event)
    case event := <-sub2.Events:
        handleDataEvent(event)
    }
}
```

---

## Testing Requirements

### Unit Tests
- Path normalization and resolution
- CRUD operations on files and directories
- Permission handling
- Error conditions (not found, already exists, etc.)

### Integration Tests
- Concurrent access scenarios
- Large directory listings
- Transaction rollback scenarios
- Volume/VFS lifecycle

### Compatibility Tests
- Drop-in replacement for os package functions
- fs.FS interface compatibility (Go 1.16+)

---

## Future Enhancements (Out of Scope for V1)
- File chunking for large files
- Compression support
- Encryption at rest
- Symbolic links
- Hard links
- Extended attributes
- File locking
- Quotas
- Snapshots/versioning
- Network access (mount over network)

---

## Performance Considerations
- Index on (volume_id, parent_id, name) for fast lookups
- Keep file content separate from metadata
- Use prepared statements
- Connection pooling
- WAL mode for better concurrency
- In-memory caching layer (optional)

---

## Security Considerations
- SQL injection prevention (use parameterized queries)
- Path traversal validation
- Permission checks before operations
- Safe handling of special characters in filenames
