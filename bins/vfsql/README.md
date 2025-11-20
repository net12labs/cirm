# VFSQL - Virtual File System in SQLite

A complete virtual file system implementation backed by SQLite, providing os-compatible file access patterns in Go.

## Features

- **OS-Compatible API**: Drop-in replacement for many `os` package functions
- **SQLite Backend**: All data stored in a single SQLite database file
- **Hierarchical Structure**: Database → Volumes → Filesystems → Directories → Files
- **Extended Metadata**: Tags and descriptions for files and directories
- **File Variants**: Store multiple versions of files (thumbnails, resized images, etc.)
- **Advanced Search**: Search by name, tags, size, time, and more
- **Event System**: Subscribe to filesystem events (create, modify, delete, etc.)
- **Transactions**: ACID guarantees for all operations
- **WAL Mode**: Better concurrency with SQLite's Write-Ahead Logging

## Installation

```bash
go get github.com/mattn/go-sqlite3
```

## Quick Start

```go
package main

import (
    "fmt"
    "io"
    "vfsql"
)

func main() {
    // Open/create database
    db, err := vfsql.Open("mydata.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create a volume
    vol, err := db.CreateVolume("myvolume")
    if err != nil {
        panic(err)
    }

    // Get the default filesystem
    vfs, err := vol.GetVFS("default")
    if err != nil {
        panic(err)
    }

    // Use like regular filesystem
    vfs.MkdirAll("/documents/reports", 0755)
    
    f, _ := vfs.Create("/documents/reports/report.txt")
    f.Write([]byte("Hello, VFSQL!"))
    f.Close()

    // Add metadata
    vfs.SetTags("/documents/reports/report.txt", []string{"important", "2024"})
    
    // Read back
    f, _ = vfs.Open("/documents/reports/report.txt")
    content, _ := io.ReadAll(f)
    f.Close()
    
    fmt.Println(string(content))
}
```

## Architecture

### Database Schema

- **volumes**: Logical partitions within the database
- **filesystems**: Named file systems within volumes
- **inodes**: File and directory metadata
- **storage**: Actual file content (separate from metadata)
- **variants**: File variant relationships

### Hierarchy

```
SQLite Database
 └── Volumes (logical partitions)
      └── Filesystems (VFS instances)
           └── Directories & Files
                └── File Variants (optional)
```

## Core Operations

### File Operations

```go
// Create/open files
f, err := vfs.Create("/file.txt")
f, err := vfs.Open("/file.txt")
f, err := vfs.OpenFile("/file.txt", os.O_RDWR, 0644)

// Read/write
n, err := f.Read(buf)
n, err := f.Write(data)
f.Seek(0, io.SeekStart)

// File operations
f.Truncate(100)
f.Sync()
f.Close()

// Delete
vfs.Remove("/file.txt")
```

### Directory Operations

```go
// Create directories
vfs.Mkdir("/dir", 0755)
vfs.MkdirAll("/path/to/dir", 0755)

// List contents
entries, err := vfs.ReadDir("/path")
for _, entry := range entries {
    fmt.Println(entry.Name(), entry.IsDir())
}

// Remove
vfs.Remove("/empty-dir")
vfs.RemoveAll("/dir-with-contents")

// Rename/move
vfs.Rename("/old/path", "/new/path")
```

### Metadata Operations

```go
// Descriptions
vfs.SetDescription("/file.txt", "My important file")
desc, err := vfs.GetDescription("/file.txt")

// Tags
vfs.SetTags("/file.txt", []string{"important", "draft"})
vfs.AddTag("/file.txt", "reviewed")
vfs.RemoveTag("/file.txt", "draft")
tags, err := vfs.GetTags("/file.txt")
```

### File Variants

```go
// Create variant
thumbnailData := resize(originalData, 200, 200)
vfs.CreateVariant("/photo.jpg", "thumb.jpg", thumbnailData)

// Access variant (two ways)
f, err := vfs.Open("/photo.jpg/vo/thumb.jpg")  // Virtual path
f, err := vfs.GetVariant("/photo.jpg", "thumb.jpg")  // Direct API

// List and remove
variants, err := vfs.ListVariants("/photo.jpg")
vfs.RemoveVariant("/photo.jpg", "thumb.jpg")
```

## Search

### Simple Search

```go
// By name pattern
paths, err := vfs.FindByName("*.pdf", &vfsql.FindOptions{
    BasePath:  "/documents",
    Recursive: true,
})

// By tags
paths, err := vfs.FindByTag([]string{"important", "2024"}, true) // AND
paths, err := vfs.FindByTag([]string{"draft", "review"}, false)  // OR
```

### Advanced Search

```go
results, err := vfs.Search(&vfsql.SearchQuery{
    // Name/pattern matching
    NamePattern: "report*.pdf",
    
    // Path constraints
    BasePath:  "/documents",
    Recursive: true,
    MaxDepth:  3,
    
    // Type and metadata
    Type:        vfsql.FileTypeFile,
    Tags:        []string{"2024", "financial"},
    TagMatchAll: true,
    Description: "quarterly",
    
    // Size constraints
    MinSize: 1024,       // 1KB
    MaxSize: 10485760,   // 10MB
    
    // Time constraints (Unix timestamps)
    ModifiedAfter: time.Now().AddDate(0, 0, -30).Unix(), // Last 30 days
    
    // Sorting and pagination
    SortBy:    vfsql.SortByModified,
    SortOrder: vfsql.Descending,
    Limit:     50,
    Offset:    0,
})

fmt.Printf("Found %d of %d total matches\n", 
    len(results.Paths), results.TotalCount)
```

## Event Subscription

```go
// Subscribe to filesystem events
sub, err := vfs.Subscribe(&vfsql.EventFilter{
    Paths:      []string{"/documents"},
    Recursive:  true,
    EventTypes: []vfsql.EventType{
        vfsql.EventCreate,
        vfsql.EventModify,
        vfsql.EventDelete,
    },
    NamePattern: "*.txt",
})
defer vfs.Unsubscribe(sub)

// Process events
go func() {
    for {
        select {
        case event := <-sub.Events:
            fmt.Printf("%v: %s\n", event.Type, event.Path)
        case err := <-sub.Errors:
            log.Printf("Error: %v\n", err)
        }
    }
}()
```

### Event Types

- `EventCreate`: File or directory created
- `EventModify`: File content modified
- `EventDelete`: File or directory deleted
- `EventRename`: File or directory renamed/moved
- `EventChmod`: Permissions changed
- `EventChown`: Ownership changed
- `EventMetadata`: Tags or description changed
- `EventVariant`: Variant created or deleted

## Implementation Details

### Path Resolution

- Paths are normalized automatically
- Supports `.` and `..` in paths
- Root directory is `/`
- Special `/vo/` virtual directory for variants

### Transaction Handling

- Read operations use shared reads
- Write operations use transactions
- File writes buffered until `Sync()` or `Close()`
- Batch operations wrapped in single transaction

### Performance Optimization

- Indexes on frequently queried fields
- Separate storage table for file content
- WAL mode for better concurrency
- Prepared statements (TODO in production)
- Connection pooling available

### Error Handling

Returns standard Go errors:
- `os.ErrNotExist`
- `os.ErrExist`
- `os.ErrPermission`
- `os.ErrInvalid`
- `fs.ErrClosed`

## Testing

```bash
go test -v
```

## Limitations (V1)

- No chunking (for small files only)
- No symbolic links
- No hard links
- No file locking
- No encryption
- No compression
- UID/GID are stored but not enforced

## Future Enhancements

- File chunking for large files
- Compression support
- Encryption at rest
- Symbolic/hard links
- File locking
- Quotas
- Snapshots/versioning
- Network access

## License

[Your License Here]
