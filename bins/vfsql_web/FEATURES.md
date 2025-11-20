# VFSQL Web Interface Features

## 🎯 Drag & Drop File Upload

### How to Use

1. **Drag and Drop**
   - Drag one or multiple files from your computer
   - Drop them into the blue drop zone in the file browser
   - Files will be uploaded to the currently selected directory

2. **Click to Browse**
   - Click anywhere in the drop zone
   - Select files using your system's file picker
   - Multiple file selection supported

3. **Upload Progress**
   - Real-time progress indicator
   - Per-file upload status (pending, success, error)
   - Cancel button to abort ongoing uploads

### Upload Target

The upload destination is shown in the drop zone:
- Default: `/` (root directory)
- Updates when you select a directory in the tree
- Updates to parent directory when you select a file

### Supported File Types

- **Text Files**: `.txt`, `.md`, `.json`, `.xml`, `.csv`, `.log`
- **Code Files**: `.js`, `.css`, `.html`, `.py`, `.go`, `.java`
- **Binary Files**: Images, documents, archives (uploaded as-is)
- **Any File Type**: All files are supported

### File Size Limits

- Maximum per file: 32 MB
- No limit on number of files per upload
- Total upload size limited by browser memory

## 📥 Download Files

### How to Download

1. Select a file in the file browser
2. Click the **📥 Download** button in the editor
3. File will be downloaded to your default downloads folder

### What Can Be Downloaded

- Any file stored in the VFS
- Original filename is preserved
- Content is downloaded exactly as stored

## 📝 File Management

### Create Files

**Method 1: New File Button**
- Click **📄 New File** button
- Enter full path (e.g., `/docs/readme.txt`)
- File is created and opened in editor

**Method 2: Upload**
- Use drag & drop or click to browse
- Files are created with their original names

### Create Folders

- Click **📁 New Folder** button
- Enter full path (e.g., `/projects/myproject`)
- Parent directories are created automatically

### Edit Files

1. Click on a file in the tree
2. Content appears in the editor
3. Make changes
4. Click **💾 Save** to persist changes
5. Auto-save indicator shows when file is modified

### Delete Files/Folders

1. Select file or folder
2. Click **🗑️ Delete** button
3. Confirm deletion
4. File/folder is removed from VFS

## 🏷️ Metadata Management

### Descriptions

- Add human-readable descriptions to files
- Searchable via the search interface
- Useful for documentation

### Tags

- Add multiple tags to files
- Tags are comma-separated keywords
- Tags are searchable
- Visual tag chips with remove buttons
- Use for organization and categorization

### File Information

- Name, size, permissions
- Created, modified, accessed times
- File type (file or directory)

## 🔄 File Variants

### What Are Variants?

Variants are alternative versions of a file stored alongside the original:
- Thumbnails for images
- Resized versions
- Converted formats
- Backups
- Processed versions

### Creating Variants

1. Select a file
2. Go to **Variants** tab
3. Enter variant name (e.g., `thumbnail.jpg`)
4. Enter variant content
5. Click **Create Variant**

### Accessing Variants

- Listed in the Variants tab
- Can be viewed directly
- Can be deleted individually
- Don't appear in normal directory listings

### Virtual Paths

Variants can be accessed via special paths:
- Original: `/photo.jpg`
- Variant: `/photo.jpg/vo/thumbnail.jpg`

## 🔍 Advanced Search

### Search Criteria

**Name Pattern**
- Use wildcards: `*.txt`, `report*.pdf`
- Case-sensitive matching
- Glob-style patterns

**Path**
- Base directory to search from
- Can be any valid path
- Defaults to `/` (root)

**Tags**
- Comma-separated tag list
- Find files with matching tags
- Useful for organization

**Description**
- Search within file descriptions
- Partial matching
- Case-insensitive

**Recursive**
- Check to search subdirectories
- Uncheck for current directory only

### Search Results

- Shows all matching files
- Click result to open file
- Shows total count
- Results are paths to files

## ⚡ Real-time Events

### Event Types

- **Create** - File or directory created
- **Modify** - File content changed
- **Delete** - File or directory removed
- **Rename** - File or directory moved
- **Chmod** - Permissions changed
- **Chown** - Ownership changed
- **Metadata** - Tags or description changed
- **Variant** - Variant created or removed

### Using Events

1. Go to **Events** tab
2. Click **▶️ Start Listening**
3. Events appear in real-time
4. Color-coded by type
5. Shows timestamp and path
6. Click **⏸️ Stop Listening** to pause
7. Click **🗑️ Clear** to clear log

### Event Log

- Scrollable log window
- Most recent events at top
- Limited to last 100 events
- Dark terminal-style display

## 🎨 User Interface

### Layout

- **Left Sidebar**: File browser and drop zone
- **Main Area**: Tabbed interface for different functions
- **Header**: Application title and description

### Navigation

- Click files/folders in tree to select
- Selected item is highlighted
- Tabs switch between different functions
- Responsive design works on mobile

### Visual Feedback

- Success messages (green)
- Error messages (red)
- Loading indicators
- Progress bars for uploads
- Active state highlighting

## 🔧 Tips & Tricks

### Bulk Operations

- Select multiple files for upload
- Drag entire folders (files extracted)
- Upload progress tracks each file

### Keyboard Shortcuts

- Browser's standard shortcuts work
- Ctrl+S in editor (triggers save button)
- Tab navigation in forms

### Performance

- Large files may take time to upload
- Upload happens in batches
- Tree refresh after operations
- Events are streamed efficiently

### Organization

- Use folders to organize files
- Tag files by project, type, status
- Add descriptions for documentation
- Create variants for different sizes

### Integration

- API can be called directly
- Events can be consumed by other apps
- Search is programmable
- All operations have REST endpoints

## 🚀 Advanced Usage

### API Integration

All web interface features are backed by REST API endpoints:
- Automate file operations
- Integrate with other systems
- Build custom interfaces
- Use in scripts and tools

### Event-Driven Workflows

- Subscribe to specific event types
- React to filesystem changes
- Build automation pipelines
- Monitor file activity

### Metadata-Based Organization

- Tag files with projects, clients, status
- Search by multiple criteria
- Build virtual collections
- Create dynamic views

## 📚 Learn More

- See README.md for API documentation
- Check examples in the test application
- Explore the vfsql package documentation
