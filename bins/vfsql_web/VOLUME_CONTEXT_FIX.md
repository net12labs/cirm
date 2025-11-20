# Volume/VFS Context Fix

## Problem
Files were always being uploaded to the default volume/VFS regardless of which volume/VFS was selected in the UI.

## Root Cause
1. The JavaScript was adding `volume` and `vfs` to FormData, but the backend `getVFS()` only checked query parameters
2. Many API calls weren't using the `getAPIUrl()` helper to include volume/VFS parameters
3. Some backend handlers (stat, tree, files, etc.) weren't calling `getVFS(r)` to get the correct VFS instance

## Changes Made

### Backend (server.go)

1. **Updated `getVFS()` to check both query params and form values:**
   ```go
   // Try query params first, then form values (for multipart uploads)
   volumeName := r.URL.Query().Get("volume")
   if volumeName == "" {
       volumeName = r.FormValue("volume")
   }
   ```

2. **Updated all handlers to use dynamic VFS:**
   - `handleTree` - ✅ Uses `getVFS(r)`
   - `handleFiles` - ✅ Uses `getVFS(r)`
   - `handleFile` - ✅ Uses `getVFS(r)`
   - `handleDirs` - ✅ Uses `getVFS(r)`
   - `handleUpload` - ✅ Uses `getVFS(r)`
   - `handleDownload` - ✅ Uses `getVFS(r)`
   - `handleStat` - ✅ Uses `getVFS(r)`

3. **Updated `buildTree()` to accept VFS parameter:**
   ```go
   func (s *Server) buildTree(vfs *vfsql.VFS, path string, depth int)
   ```

### Frontend (app.js)

1. **Updated all API calls to use `getAPIUrl()`:**
   - `loadTextFile()` - ✅
   - `loadImageFile()` - ✅
   - `loadBinaryFile()` - ✅
   - `saveFile()` - ✅
   - `deleteFile()` - ✅
   - `createFile()` - ✅
   - `createFolder()` - ✅
   - `refreshTree()` - ✅

2. **getAPIUrl() automatically adds volume/VFS params:**
   ```javascript
   function getAPIUrl(endpoint, extraParams = {}) {
       const url = new URL(`${API_BASE}${endpoint}`, window.location.origin);
       if (currentVolume) url.searchParams.set('volume', currentVolume);
       if (currentVFS) url.searchParams.set('vfs', currentVFS);
       // ...
   }
   ```

## Testing

1. **Build and run:**
   ```bash
   cd /home/lxk/Desktop/cirm/bins/vfsql_web
   make clean && make
   ./vfsql_web
   ```

2. **Test volume switching:**
   - Create a new volume
   - Create a new VFS in that volume
   - Click on the VFS to select it
   - Upload files - they should go to the selected VFS
   - File tree should refresh to show files in that VFS

3. **Test all operations:**
   - ✅ Upload files
   - ✅ Create files/folders
   - ✅ View images
   - ✅ Edit text files
   - ✅ Download files
   - ✅ Delete files
   - ✅ Search files
   - ✅ Get CDN URLs

## Result

All file operations now respect the currently selected volume and VFS. Users can:
- Switch between volumes and VFS instances
- Upload files to the correct location
- View and manage files in different VFS instances
- All operations are isolated per volume/VFS

✅ **Issue resolved!**
