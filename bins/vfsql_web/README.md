# VFSQL Web Interface

A complete web-based interface for testing and using the VFSQL virtual file system.

## Features

- 📁 **File Browser** - Navigate and manage files and directories
- ✏️ **Built-in Editor** - Create and edit files directly in the browser
- 🏷️ **Metadata Management** - Add descriptions and tags to files
- 🔄 **File Variants** - Create and manage file variants
- 🔍 **Advanced Search** - Search by name, tags, descriptions, and more
- ⚡ **Real-time Events** - Monitor filesystem changes in real-time
- 🎨 **Modern UI** - Clean, responsive interface

## Quick Start

### Method 1: Using Make (Recommended)

```bash
cd /home/lxk/Desktop/cirm/bins/vfsql_web

# Build and run (automatically sets up web directory)
make
```

### Method 2: Using the Run Script

```bash
cd /home/lxk/Desktop/cirm/bins/vfsql_web

# Make script executable
chmod +x run.sh

# Run
./run.sh
```

### Method 3: Manual Build

```bash
cd /home/lxk/Desktop/cirm/bins/vfsql_web

# Create web symlink
ln -sf ../vfsql/web web

# Build
go build -o vfsql_web .

# Run
./vfsql_web
```

### Options

```bash
./vfsql_web -db mydata.db -port 8080
```

- `-db` - Path to database file (default: `vfsql.db`)
- `-port` - HTTP server port (default: `8080`)

### Access

Open your browser and navigate to:
- **Web Interface**: http://localhost:8080
- **API Docs**: http://localhost:8080/api

## Web Interface

### File Browser Tab

- **New File** - Create a new file
- **New Folder** - Create a new directory
- **Refresh** - Reload the file tree
- **Drag & Drop** - Drag files from your computer directly into the drop zone
- **Click to Browse** - Click the drop zone to select files from your computer
- Click on any file to view/edit its contents
- Upload target path updates based on selected directory

### Editor Tab

- View and edit file contents
- **Upload Files** button - Select files from your computer
- **Download** button - Download current file to your computer
- **Save** button - Persist changes to the file
- **Delete** button - Remove files/directories
- Auto-save indicator when changes are made
- Upload progress indicator with per-file status

### Metadata Tab

- Add descriptions to files
- Manage tags (add/remove)
- View file information (size, mode, modified time)

### Variants Tab

- Create variants of files (e.g., thumbnails, resized images)
- View existing variants
- Delete variants

### Search Tab

- Search by file pattern (*.txt, *.jpg, etc.)
- Search by tags
- Search in descriptions
- Recursive or non-recursive search
- Click results to open files

### Events Tab

- Real-time monitoring of filesystem events
- See creates, modifications, deletions, etc.
- Event log with timestamps
- Start/stop listening to events

## API Endpoints

### Files

- `GET /api/files` - List files in root
- `POST /api/files` - Create new file
- `GET /api/file/:path` - Get file content
- `PUT /api/file/:path` - Update file content
- `DELETE /api/file/:path` - Delete file

### Directories

- `POST /api/dirs` - Create directory
- `GET /api/dir/:path` - List directory contents
- `DELETE /api/dir/:path` - Delete directory

### Metadata

- `GET /api/metadata/:path` - Get metadata
- `PUT /api/metadata/:path` - Set metadata
- `POST /api/tags/:path` - Add tag
- `DELETE /api/tags/:path?tag=:tag` - Remove tag

### Variants

- `GET /api/variants/:path` - List variants
- `POST /api/variants/:path` - Create variant
- `DELETE /api/variants/:path?name=:name` - Delete variant

### Search

- `GET /api/search?pattern=:pattern&path=:path&tags=:tags&recursive=:bool` - Search files

### Utilities

- `GET /api/stat/:path` - Get file stats
- `GET /api/tree?path=:path` - Get directory tree
- `GET /api/events` - Server-sent events stream
- `POST /api/upload` - Multipart file upload (supports binary files)

## Example Usage

### Using cURL

```bash
# Create a file
curl -X POST http://localhost:8080/api/files \
  -H "Content-Type: application/json" \
  -d '{"path":"/test.txt","content":"Hello, World!"}'

# Read a file
curl http://localhost:8080/api/file/test.txt

# Add tags
curl -X PUT http://localhost:8080/api/metadata/test.txt \
  -H "Content-Type: application/json" \
  -d '{"description":"Test file","tags":["test","example"]}'

# Search
curl "http://localhost:8080/api/search?tags=test&recursive=true"

# Create variant
curl -X POST http://localhost:8080/api/variants/test.txt \
  -H "Content-Type: application/json" \
  -d '{"name":"backup.txt","content":"Backup version"}'

# Upload files (multipart)
curl -X POST http://localhost:8080/api/upload \
  -F "path=/uploads" \
  -F "files=@/path/to/file1.txt" \
  -F "files=@/path/to/file2.jpg"
```

### Using JavaScript

```javascript
// Create file
await fetch('http://localhost:8080/api/files', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        path: '/hello.txt',
        content: 'Hello from JavaScript!'
    })
});

// Search
const results = await fetch(
    'http://localhost:8080/api/search?pattern=*.txt&recursive=true'
).then(r => r.json());

console.log('Found files:', results.paths);
```

## Architecture

```
vfsql_web/
├── main.go           # Server entry point
├── go.mod            # Go module
└── ../vfsql/
    ├── api/
    │   └── server.go # HTTP API handlers
    └── web/
        ├── index.html # Web interface
        ├── style.css  # Styles
        └── app.js     # JavaScript app
```

## Development

### Project Structure

The web server is built on top of the VFSQL package:

1. **API Layer** (`api/server.go`) - REST API endpoints
2. **Web Layer** (`web/`) - Static HTML/CSS/JS files
3. **VFSQL Core** - File system implementation

### Adding New Features

1. Add API endpoint in `api/server.go`
2. Add UI controls in `web/index.html`
3. Add styling in `web/style.css`
4. Add JavaScript logic in `web/app.js`

## Troubleshooting

### 404 Page Not Found

If you're getting 404 errors:

1. **Check web directory symlink:**
   ```bash
   ls -la web
   # Should show: web -> ../vfsql/web
   ```

2. **Create symlink manually:**
   ```bash
   ln -sf ../vfsql/web web
   ```

3. **Verify web files exist:**
   ```bash
   ls -la ../vfsql/web/
   # Should show: index.html, style.css, app.js
   ```

4. **Run from correct directory:**
   ```bash
   # Make sure you're in vfsql_web directory
   pwd
   # Should show: .../cirm/bins/vfsql_web
   ```

### Port Already in Use

```bash
# Use a different port
./vfsql_web -port 8081
```

### Database Locked

Make sure no other process is using the database file.

### CORS Issues

The server enables CORS by default. If you still have issues, check your browser console.

### Web Files Not Found

If the web interface shows a fallback page:
- The API is working correctly
- Web files couldn't be found
- Use the API directly or fix the web directory path

## Security Note

This is a development/testing tool. For production use:
- Add authentication
- Add HTTPS
- Validate all inputs
- Add rate limiting
- Restrict file access

## License

Same as parent project.
