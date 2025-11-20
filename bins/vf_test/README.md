# VFSQL Test Application

This is a test application demonstrating the use of the VFSQL virtual file system as a module.

## Overview

This application imports and uses the `vfsql` package from `../vfsql` to demonstrate:

- Basic file operations (create, read, write)
- Directory operations (mkdir, readdir, etc.)
- Metadata management (tags and descriptions)
- File variants
- Advanced search capabilities
- Event subscription

## Building and Running

### Using Make (recommended)

```bash
# Build and run
make

# Or individually
make build
make run

# Clean up
make clean
```

### Using Go directly

```bash
# Build
go build -o vf_test .

# Run
./vf_test

# Or in one step
go run main.go
```

## What the Test Does

The test application creates a database called `test.db` and performs the following operations:

1. **Basic File Operations**
   - Creates and writes to files
   - Reads file content back
   - Gets file information

2. **Directory Operations**
   - Creates nested directory structures
   - Creates multiple files
   - Lists directory contents

3. **Metadata Operations**
   - Sets descriptions on files
   - Adds and manages tags
   - Retrieves metadata

4. **File Variants**
   - Creates an original file
   - Creates multiple variants (thumbnail, resized)
   - Accesses variants via virtual paths
   - Lists all variants

5. **Search Operations**
   - Simple name pattern searches
   - Tag-based searches
   - Advanced searches with multiple criteria

6. **Event Subscription**
   - Subscribes to filesystem events
   - Creates, modifies, and deletes files
   - Captures and displays events in real-time

## Output

The application produces detailed output showing each operation and its results. After completion, you'll have a `test.db` file that you can inspect:

```bash
# Inspect the database with sqlite3
sqlite3 test.db

# Try these queries:
sqlite> .tables
sqlite> SELECT * FROM volumes;
sqlite> SELECT name, type, size FROM inodes LIMIT 10;
sqlite> SELECT variant_name FROM variants;
```

## Module Usage

This project demonstrates proper module usage with a local replace directive in `go.mod`:

```go
module github.com/lxk/cirm/bins/vf_test

go 1.21

require (
    github.com/lxk/cirm/bins/vfsql v0.0.0
    github.com/mattn/go-sqlite3 v1.14.18
)

replace github.com/lxk/cirm/bins/vfsql => ../vfsql
```

The `replace` directive allows us to use the local version of vfsql during development.

## Dependencies

- Go 1.21 or later
- github.com/mattn/go-sqlite3
- github.com/lxk/cirm/bins/vfsql (local module)

## Project Structure

```
vf_test/
├── go.mod              # Go module definition
├── main.go             # Test application
├── Makefile            # Build automation
└── README.md           # This file
```

## Troubleshooting

### CGO_ENABLED error

If you get a CGO error when building, ensure you have a C compiler installed:

```bash
# On Ubuntu/Debian
sudo apt-get install build-essential

# On macOS
xcode-select --install
```

### Module not found

If you get a module not found error, make sure you're in the correct directory and run:

```bash
go mod tidy
```

## Next Steps

After running the test application, you can:

1. Examine the generated `test.db` file
2. Modify `main.go` to test additional features
3. Create your own application using vfsql
4. Explore the vfsql package documentation in `../vfsql/README.md`

## License

Same as parent project.
