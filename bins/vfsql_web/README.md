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

- **New File** - Create a new file in the currently selected directory
- **New Folder** - Create a new directory in the currently selected directory
- **Refresh** - Reload the file tree
- **Drag & Drop** - Drag files from your computer directly into the drop zone
- **Click to Browse** - Click the drop zone to select files from your computer
- Click on any file to view/edit its contents
- Click on any directory to set it as the target for new files/folders/uploads
- Current target directory is shown in the drop zone ("Uploading to: /path")

### Editor Tab

- **Text Files** - View and edit in the built-in editor
- **Image Files** - Automatic image preview with dimensions and size
  - Supported: JPG, PNG, GIF, BMP, WEBP, SVG, ICO
- **Binary Files** - Shows file information (PDF, ZIP, etc.)
- **Upload Files** button - Select files from your computer
- **Download** button - Download current file to your computer
- **Save** button - Persist changes to text files
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
- `GET /api/cdn-url/:path` - Generate CDN URL for a file

### CDN Endpoint

- `GET /cdn/{hash}/{path}` - Serve files with CDN caching headers
- `HEAD /cdn/{hash}/{path}` - Get file headers without content

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

## CDN URLs

VFSQL provides CDN-style URLs for serving files that can be embedded in other websites.

### Features

- **Long-term caching** - Files are cached for 1 year with immutable directive
- **CORS enabled** - Can be used from any domain
- **ETag support** - Efficient caching with 304 Not Modified responses
- **Multiple content types** - Automatic detection for images, videos, fonts, etc.
- **High performance** - Optimized for static asset delivery

### Getting CDN URLs

#### Via Web Interface

1. Select a file in the file browser
2. Click the **🔗 Get CDN URL** button
3. Copy the generated URL
4. Use it in your HTML, CSS, or JavaScript

#### Via API

```bash
# Get CDN URL for a file
curl http://localhost:8080/api/cdn-url/images/logo.png

# Response:
{
  "url": "/cdn/a1b2c3d4e5f6/images/logo.png",
  "full_url": "http://localhost:8080/cdn/a1b2c3d4e5f6/images/logo.png",
  "path": "/images/logo.png",
  "hash": "a1b2c3d4e5f6",
  "size": 12345,
  "modified": 1234567890
}
```

### Using CDN URLs

#### In HTML

```html
<!-- Images -->
<img src="http://your-server.com/cdn/a1b2c3d4e5f6/images/photo.jpg" alt="Photo">

<!-- CSS -->
<link rel="stylesheet" href="http://your-server.com/cdn/a1b2c3d4e5f6/styles/main.css">

<!-- JavaScript -->
<script src="http://your-server.com/cdn/a1b2c3d4e5f6/scripts/app.js"></script>

<!-- Fonts -->
@font-face {
  font-family: 'MyFont';
  src: url('http://your-server.com/cdn/a1b2c3d4e5f6/fonts/myfont.woff2');
}

<!-- Videos -->
<video src="http://your-server.com/cdn/a1b2c3d4e5f6/videos/intro.mp4" controls></video>

<!-- PDFs -->
<embed src="http://your-server.com/cdn/a1b2c3d4e5f6/docs/manual.pdf" type="application/pdf">
```

#### In Markdown

```markdown
![Logo](http://your-server.com/cdn/a1b2c3d4e5f6/images/logo.png)
[Download PDF](http://your-server.com/cdn/a1b2c3d4e5f6/docs/guide.pdf)
```

### URL Format

```
/cdn/{hash}/{path}
```

- **hash** - Security hash (prevents path enumeration)
- **path** - File path within the VFS

### Supported File Types

The CDN endpoint automatically sets correct Content-Type headers for:

- **Images**: JPG, PNG, GIF, SVG, WEBP, ICO
- **Documents**: PDF, TXT, JSON, XML, HTML, MD, CSV
- **Styles**: CSS
- **Scripts**: JS
- **Fonts**: WOFF, WOFF2, TTF, OTF
- **Media**: MP4, MP3
- **Archives**: ZIP

### Cache Headers

```
Cache-Control: public, max-age=31536000, immutable
Access-Control-Allow-Origin: *
X-Content-Type-Options: nosniff
ETag: "timestamp-size"
```

### Example Use Cases

1. **Serve images for a blog** - Host images in VFSQL, link from blog posts
2. **CDN for web assets** - Serve CSS, JS, fonts for websites
3. **File sharing** - Share files via permanent URLs
4. **API documentation** - Host images and diagrams
5. **Email templates** - Link to hosted images in emails
6. **Mobile apps** - Serve assets for mobile applications
