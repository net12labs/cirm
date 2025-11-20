// filepath: /home/lxk/Desktop/cirm/bins/vfsql/web/app.js
// API base URL
const API_BASE = 'http://localhost:8080/api';

// State
let currentPath = null;
let currentIsDir = false;
let eventsSource = null;
let isListeningToEvents = false;
let uploadTargetPath = '/';
let uploadAbortController = null;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    refreshTree();
    document.getElementById('editor').addEventListener('input', () => {
        document.getElementById('saveBtn').disabled = false;
    });
    setupDragAndDrop();
    
    // Set initial upload path display
    document.getElementById('uploadPath').textContent = '/';
});

// Drag and Drop Setup
function setupDragAndDrop() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    
    // Click to browse
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });
    
    // File input change
    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length > 0) {
            handleFiles(e.target.files);
        }
    });
    
    // Prevent default drag behaviors
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
        document.body.addEventListener(eventName, preventDefaults, false);
    });
    
    // Highlight drop zone when item is dragged over it
    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.add('drag-over');
        }, false);
    });
    
    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => {
            dropZone.classList.remove('drag-over');
        }, false);
    });
    
    // Handle dropped files
    dropZone.addEventListener('drop', (e) => {
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFiles(files);
        }
    }, false);
}

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

// Handle file uploads
async function handleFiles(files) {
    const fileArray = Array.from(files);
    
    if (fileArray.length === 0) return;
    
    // Show progress UI
    const progressDiv = document.getElementById('uploadProgress');
    const progressFill = document.getElementById('progressFill');
    const progressText = document.getElementById('progressText');
    const detailsDiv = document.getElementById('uploadDetails');
    
    progressDiv.classList.remove('hidden');
    progressText.textContent = `Uploading ${fileArray.length} file(s)...`;
    detailsDiv.innerHTML = '';
    
    uploadAbortController = new AbortController();
    
    // Use multipart upload for better binary support
    const formData = new FormData();
    
    // Ensure path is set correctly
    const targetPath = uploadTargetPath || '/';
    formData.append('path', targetPath);
    
    console.log(`Uploading ${fileArray.length} file(s) to: ${targetPath}`);
    
    // Add all files to form data
    fileArray.forEach(file => {
        formData.append('files', file);
        
        const itemDiv = document.createElement('div');
        itemDiv.className = 'upload-item';
        itemDiv.id = `upload-${file.name}`;
        itemDiv.innerHTML = `
            <span class="name">${file.name} → ${targetPath}</span>
            <span class="status pending">⏳</span>
        `;
        detailsDiv.appendChild(itemDiv);
    });
    
    try {
        // Upload all files at once
        console.log('Sending upload request...');
        const response = await fetch(`${API_BASE}/upload`, {
            method: 'POST',
            body: formData,
            signal: uploadAbortController.signal
        });
        
        console.log('Upload response status:', response.status);
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Upload failed: ${response.statusText} - ${errorText}`);
        }
        
        const result = await response.json();
        console.log('Upload result:', result);
        
        // Update status for each file
        if (result.results && Array.isArray(result.results)) {
            result.results.forEach(fileResult => {
                const itemDiv = document.getElementById(`upload-${fileResult.name}`);
                if (itemDiv) {
                    if (fileResult.success) {
                        itemDiv.querySelector('.status').textContent = '✓';
                        itemDiv.querySelector('.status').className = 'status success';
                        // Update to show actual path
                        if (fileResult.path) {
                            itemDiv.querySelector('.name').textContent = fileResult.path;
                        }
                    } else {
                        itemDiv.querySelector('.status').textContent = '✗';
                        itemDiv.querySelector('.status').className = 'status error';
                        itemDiv.title = fileResult.error || 'Upload failed';
                    }
                }
            });
            
            progressFill.style.width = '100%';
            
            const successful = result.results.filter(r => r.success).length;
            progressText.textContent = `Uploaded ${successful} of ${fileArray.length} file(s)`;
            
            if (successful > 0) {
                showSuccess(`Successfully uploaded ${successful} file(s) to ${targetPath}`);
            }
            if (successful < fileArray.length) {
                showError(`Failed to upload ${fileArray.length - successful} file(s)`);
            }
        } else {
            throw new Error('Invalid response format from server');
        }
        
    } catch (error) {
        if (error.name === 'AbortError') {
            progressText.textContent = 'Upload cancelled';
        } else {
            progressText.textContent = 'Upload failed';
            showError('Upload error: ' + error.message);
        }
        
        // Mark all pending as error
        detailsDiv.querySelectorAll('.status.pending').forEach(status => {
            status.textContent = '✗';
            status.className = 'status error';
        });
    }
    
    // Refresh tree and hide progress after delay
    refreshTree();
    setTimeout(() => {
        progressDiv.classList.add('hidden');
        uploadAbortController = null;
    }, 3000);
}

// Note: This function is no longer used since we switched to multipart upload
function readFileContent(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (e) => resolve(e.target.result);
        reader.onerror = (e) => reject(e);
        
        // Read as text for now (could add binary support later)
        if (file.type.startsWith('text/') || file.name.match(/\.(txt|md|json|xml|csv|log|js|css|html)$/i)) {
            reader.readAsText(file);
        } else {
            // For binary files, read as base64
            reader.readAsDataURL(file);
        }
    });
}

function cancelUpload() {
    if (uploadAbortController) {
        uploadAbortController.abort();
        showError('Upload cancelled');
    }
}

function uploadFiles() {
    document.getElementById('fileInput').click();
}

// Get CDN URL for file
async function getCDNUrl() {
    if (!currentPath || currentIsDir) return;
    
    try {
        const response = await fetch(`${API_BASE}/cdn-url${currentPath}`);
        
        if (!response.ok) {
            throw new Error(`Failed to get CDN URL: ${response.statusText}`);
        }
        
        const data = await response.json();
        
        // Create a modal/dialog to show the URL
        const modal = document.createElement('div');
        modal.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.5);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 10000;
        `;
        
        modal.innerHTML = `
            <div style="background: white; padding: 30px; border-radius: 10px; max-width: 600px; width: 90%;">
                <h3 style="margin-top: 0; color: #667eea;">🔗 CDN URL</h3>
                <p style="color: #666;">Use this URL to link to your file from other websites:</p>
                
                <div style="margin: 20px 0;">
                    <label style="display: block; font-weight: 500; margin-bottom: 5px;">Full URL:</label>
                    <input type="text" id="cdnFullUrl" value="${data.full_url}" 
                           style="width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 5px; font-family: monospace;"
                           readonly>
                </div>
                
                <div style="margin: 20px 0;">
                    <label style="display: block; font-weight: 500; margin-bottom: 5px;">Relative URL:</label>
                    <input type="text" id="cdnRelUrl" value="${data.url}" 
                           style="width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 5px; font-family: monospace;"
                           readonly>
                </div>
                
                <div style="background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0;">
                    <strong>Features:</strong>
                    <ul style="margin: 10px 0; padding-left: 20px;">
                        <li>✓ Long-term caching (1 year)</li>
                        <li>✓ CORS enabled for cross-origin use</li>
                        <li>✓ ETag support for efficient caching</li>
                        <li>✓ Immutable content</li>
                    </ul>
                </div>
                
                <div style="display: flex; gap: 10px; justify-content: flex-end;">
                    <button onclick="copyToClipboard('${data.full_url}')" 
                            style="padding: 10px 20px; background: #667eea; color: white; border: none; border-radius: 5px; cursor: pointer;">
                        📋 Copy Full URL
                    </button>
                    <button onclick="this.closest('[style*=fixed]').remove()" 
                            style="padding: 10px 20px; background: #666; color: white; border: none; border-radius: 5px; cursor: pointer;">
                        Close
                    </button>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        
        // Select the URL text on click
        document.getElementById('cdnFullUrl').addEventListener('click', function() {
            this.select();
        });
        document.getElementById('cdnRelUrl').addEventListener('click', function() {
            this.select();
        });
        
        // Close on background click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
        
        showSuccess('CDN URL generated');
    } catch (error) {
        showError('Failed to get CDN URL: ' + error.message);
        console.error('CDN URL error:', error);
    }
}

// Copy to clipboard helper
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showSuccess('URL copied to clipboard!');
    }).catch(err => {
        showError('Failed to copy: ' + err.message);
    });
}

// Download file
async function downloadFile() {
    if (!currentPath || currentIsDir) return;
    
    try {
        // Use the download endpoint which sets proper headers
        const downloadUrl = `${API_BASE}/download${currentPath}`;
        
        // Fetch the file as a blob
        const response = await fetch(downloadUrl);
        
        if (!response.ok) {
            throw new Error(`Download failed: ${response.statusText}`);
        }
        
        // Get the blob
        const blob = await response.blob();
        
        console.log('Downloaded blob size:', blob.size, 'bytes');
        
        if (blob.size === 0) {
            throw new Error('Downloaded file is empty');
        }
        
        // Create object URL and download
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = currentPath.split('/').pop();
        document.body.appendChild(a);
        a.click();
        
        // Cleanup
        setTimeout(() => {
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        }, 100);
        
        showSuccess(`File downloaded (${formatSize(blob.size)})`);
    } catch (error) {
        showError('Failed to download file: ' + error.message);
        console.error('Download error:', error);
    }
}

// Tab switching
function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(`${tabName}-tab`).classList.add('active');
    event.target.classList.add('active');
    
    // Load data for specific tabs
    if (tabName === 'metadata' && currentPath) {
        loadMetadata();
    } else if (tabName === 'variants' && currentPath && !currentIsDir) {
        loadVariants();
    }
}

// File tree
async function refreshTree() {
    try {
        const tree = await fetch(`${API_BASE}/tree?path=/`).then(r => r.json());
        renderTree(tree);
    } catch (error) {
        showError('Failed to load file tree: ' + error.message);
    }
}

function renderTree(node, container = document.getElementById('fileTree'), level = 0) {
    if (!container) return;
    if (level === 0) container.innerHTML = '';
    
    const item = document.createElement('div');
    item.className = `tree-item ${node.isDir ? 'dir' : 'file'}`;
    item.textContent = node.name || '/';
    item.style.marginLeft = `${level * 15}px`;
    item.onclick = () => selectFile(node.path, node.isDir);
    container.appendChild(item);
    
    if (node.children && node.children.length > 0) {
        const childContainer = document.createElement('div');
        childContainer.className = 'tree-children';
        container.appendChild(childContainer);
        
        node.children.forEach(child => renderTree(child, childContainer, level + 1));
    }
}

// File selection
async function selectFile(path, isDir) {
    // Normalize path - remove double slashes
    path = path.replace(/\/+/g, '/');
    currentPath = path;
    currentIsDir = isDir;
    
    // Update upload target path
    if (isDir) {
        uploadTargetPath = path;
        document.getElementById('uploadPath').textContent = path;
    } else {
        // Use parent directory
        uploadTargetPath = path.substring(0, path.lastIndexOf('/')) || '/';
        document.getElementById('uploadPath').textContent = uploadTargetPath;
    }
    
    console.log(`Selected: ${path} (isDir: ${isDir}), upload target: ${uploadTargetPath}`);
    
    // Update active state
    document.querySelectorAll('.tree-item').forEach(item => {
        item.classList.remove('active');
    });
    event.target.classList.add('active');
    
    document.getElementById('currentFile').textContent = path;
    document.getElementById('saveBtn').disabled = true;
    document.getElementById('deleteBtn').disabled = false;
    document.getElementById('downloadBtn').disabled = isDir;
    document.getElementById('cdnBtn').disabled = isDir;
    
    if (!isDir) {
        await loadFile(path);
    } else {
        document.getElementById('editor').value = `[Directory: ${path}]`;
        document.getElementById('editor').disabled = true;
        document.getElementById('saveBtn').disabled = true;
    }
}

// Load file content
async function loadFile(path) {
    try {
        // Hide all content views
        document.getElementById('editor').classList.remove('hidden');
        document.getElementById('imagePreview').classList.add('hidden');
        document.getElementById('binaryInfo').classList.add('hidden');
        
        // Check if it's an image file
        const ext = path.split('.').pop().toLowerCase();
        const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg', 'ico'];
        const binaryExts = ['pdf', 'zip', 'tar', 'gz', 'exe', 'dll', 'so', 'bin'];
        
        if (imageExts.includes(ext)) {
            await loadImageFile(path);
        } else if (binaryExts.includes(ext)) {
            await loadBinaryFile(path);
        } else {
            await loadTextFile(path);
        }
    } catch (error) {
        showError('Failed to load file: ' + error.message);
    }
}

// Load text file
async function loadTextFile(path) {
    const data = await fetch(`${API_BASE}/file${path}`).then(r => r.json());
    document.getElementById('editor').value = data.content || '';
    document.getElementById('editor').disabled = false;
    document.getElementById('editor').classList.remove('hidden');
    document.getElementById('saveBtn').disabled = true;
}

// Load image file
async function loadImageFile(path) {
    // Hide editor, show image preview
    document.getElementById('editor').classList.add('hidden');
    document.getElementById('imagePreview').classList.remove('hidden');
    document.getElementById('saveBtn').disabled = true;
    
    // Use download endpoint to get the image
    const imageUrl = `${API_BASE}/download${path}`;
    const img = document.getElementById('previewImage');
    
    img.onload = function() {
        document.getElementById('imageDimensions').textContent = 
            `${img.naturalWidth} × ${img.naturalHeight} pixels`;
    };
    
    img.src = imageUrl;
    
    // Get file size
    try {
        const info = await fetch(`${API_BASE}/stat${path}`).then(r => r.json());
        document.getElementById('imageSize').textContent = formatSize(info.size);
    } catch (e) {
        document.getElementById('imageSize').textContent = '';
    }
}

// Load binary file
async function loadBinaryFile(path) {
    // Hide editor, show binary info
    document.getElementById('editor').classList.add('hidden');
    document.getElementById('binaryInfo').classList.remove('hidden');
    document.getElementById('saveBtn').disabled = true;
    
    // Get file info
    try {
        const info = await fetch(`${API_BASE}/stat${path}`).then(r => r.json());
        const filename = path.split('/').pop();
        const ext = filename.split('.').pop().toUpperCase();
        
        document.getElementById('binaryDetails').innerHTML = `
            <div class="binary-detail-row">
                <span class="binary-detail-label">Filename:</span>
                <span class="binary-detail-value">${filename}</span>
            </div>
            <div class="binary-detail-row">
                <span class="binary-detail-label">Type:</span>
                <span class="binary-detail-value">${ext} File</span>
            </div>
            <div class="binary-detail-row">
                <span class="binary-detail-label">Size:</span>
                <span class="binary-detail-value">${formatSize(info.size)}</span>
            </div>
            <div class="binary-detail-row">
                <span class="binary-detail-label">Modified:</span>
                <span class="binary-detail-value">${new Date(info.modTime).toLocaleString()}</span>
            </div>
        `;
    } catch (e) {
        console.error('Failed to load binary file info:', e);
    }
}

// Save file
async function saveFile() {
    if (!currentPath || currentIsDir) return;
    
    // Only save if editor is visible (text files)
    if (document.getElementById('editor').classList.contains('hidden')) {
        showError('Cannot save binary files');
        return;
    }
    
    try {
        const content = document.getElementById('editor').value;
        await fetch(`${API_BASE}/file${currentPath}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ content })
        });
        
        document.getElementById('saveBtn').disabled = true;
        showSuccess('File saved successfully');
    } catch (error) {
        showError('Failed to save file: ' + error.message);
    }
}

// Delete file/directory
async function deleteFile() {
    if (!currentPath) return;
    
    if (!confirm(`Are you sure you want to delete ${currentPath}?`)) return;
    
    try {
        if (currentIsDir) {
            await fetch(`${API_BASE}/dir${currentPath}`, { method: 'DELETE' });
        } else {
            await fetch(`${API_BASE}/file${currentPath}`, { method: 'DELETE' });
        }
        
        currentPath = null;
        document.getElementById('editor').value = '';
        document.getElementById('currentFile').textContent = 'No file selected';
        document.getElementById('deleteBtn').disabled = true;
        showSuccess('Deleted successfully');
        refreshTree();
    } catch (error) {
        showError('Failed to delete: ' + error.message);
    }
}

// Create new file
async function createFile() {
    // Use current upload target path (which is the selected directory)
    const basePath = uploadTargetPath || '/';
    const displayPath = basePath === '/' ? 'root' : basePath;
    const fileName = prompt(`Create new file in: ${displayPath}\n\nEnter filename (e.g., myfile.txt):`);
    if (!fileName) return;
    
    // Clean up the file name (remove leading/trailing slashes)
    const cleanName = fileName.replace(/^\/+|\/+$/g, '');
    
    if (!cleanName) {
        showError('Invalid filename');
        return;
    }
    
    // Build the full path
    let fullPath;
    if (basePath === '/') {
        fullPath = `/${cleanName}`;
    } else {
        fullPath = `${basePath}/${cleanName}`;
    }
    
    try {
        await fetch(`${API_BASE}/files`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ path: fullPath, content: '' })
        });
        
        showSuccess(`File created: ${fullPath}`);
        refreshTree();
        selectFile(fullPath, false);
    } catch (error) {
        showError('Failed to create file: ' + error.message);
    }
}

// Create new folder
async function createFolder() {
    // Use current upload target path (which is the selected directory)
    const basePath = uploadTargetPath || '/';
    const displayPath = basePath === '/' ? 'root' : basePath;
    const folderName = prompt(`Create new folder in: ${displayPath}\n\nEnter folder name (e.g., documents):`);
    if (!folderName) return;
    
    // Clean up the folder name (remove leading/trailing slashes)
    const cleanName = folderName.replace(/^\/+|\/+$/g, '');
    
    if (!cleanName) {
        showError('Invalid folder name');
        return;
    }
    
    // Build the full path
    let fullPath;
    if (basePath === '/') {
        fullPath = `/${cleanName}`;
    } else {
        fullPath = `${basePath}/${cleanName}`;
    }
    
    try {
        await fetch(`${API_BASE}/dirs`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ path: fullPath })
        });
        
        showSuccess(`Folder created: ${fullPath}`);
        refreshTree();
    } catch (error) {
        showError('Failed to create folder: ' + error.message);
    }
}

// Test file creation and download
async function testCreate() {
    try {
        console.log('Creating test file...');
        const response = await fetch(`${API_BASE}/test-create`, {
            method: 'POST'
        });
        
        const result = await response.json();
        console.log('Test result:', result);
        
        if (result.content_match) {
            showSuccess(`Test passed! File created with ${result.written_bytes} bytes, read back ${result.read_back_len} bytes`);
        } else {
            showError('Test failed: content mismatch');
        }
        
        refreshTree();
        
        // Try to download the test file
        setTimeout(() => {
            selectFile('/test-file.txt', false);
        }, 500);
        
    } catch (error) {
        showError('Test failed: ' + error.message);
        console.error('Test error:', error);
    }
}

// Metadata
async function loadMetadata() {
    if (!currentPath) return;
    
    try {
        const data = await fetch(`${API_BASE}/metadata${currentPath}`).then(r => r.json());
        document.getElementById('metaDescription').value = data.description || '';
        document.getElementById('metaTags').value = (data.tags || []).join(', ');
        
        // Display current tags
        const tagsContainer = document.getElementById('currentTags');
        tagsContainer.innerHTML = '';
        (data.tags || []).forEach(tag => {
            const tagEl = document.createElement('span');
            tagEl.className = 'tag';
            tagEl.innerHTML = `${tag} <span class="remove" onclick="removeTag('${tag}')">×</span>`;
            tagsContainer.appendChild(tagEl);
        });
        
        // Load file info
        const info = await fetch(`${API_BASE}/stat${currentPath}`).then(r => r.json());
        const infoHtml = `
            <div class="info-row"><span class="info-label">Name:</span><span class="info-value">${info.name}</span></div>
            <div class="info-row"><span class="info-label">Size:</span><span class="info-value">${formatSize(info.size)}</span></div>
            <div class="info-row"><span class="info-label">Mode:</span><span class="info-value">${info.mode}</span></div>
            <div class="info-row"><span class="info-label">Modified:</span><span class="info-value">${new Date(info.modTime).toLocaleString()}</span></div>
            <div class="info-row"><span class="info-label">Type:</span><span class="info-value">${info.isDir ? 'Directory' : 'File'}</span></div>
        `;
        document.getElementById('fileInfo').innerHTML = infoHtml;
    } catch (error) {
        showError('Failed to load metadata: ' + error.message);
    }
}

async function saveMetadata() {
    if (!currentPath) return;
    
    try {
        const description = document.getElementById('metaDescription').value;
        const tagsStr = document.getElementById('metaTags').value;
        const tags = tagsStr.split(',').map(t => t.trim()).filter(t => t);
        
        await fetch(`${API_BASE}/metadata${currentPath}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ description, tags })
        });
        
        showSuccess('Metadata saved successfully');
        loadMetadata();
    } catch (error) {
        showError('Failed to save metadata: ' + error.message);
    }
}

async function removeTag(tag) {
    if (!currentPath) return;
    
    try {
        await fetch(`${API_BASE}/tags${currentPath}?tag=${encodeURIComponent(tag)}`, {
            method: 'DELETE'
        });
        
        showSuccess('Tag removed');
        loadMetadata();
    } catch (error) {
        showError('Failed to remove tag: ' + error.message);
    }
}

// Variants
async function loadVariants() {
    if (!currentPath || currentIsDir) return;
    
    try {
        const variants = await fetch(`${API_BASE}/variants${currentPath}`).then(r => r.json());
        const listEl = document.getElementById('variantsList');
        
        if (variants.length === 0) {
            listEl.innerHTML = '<li style="color: #999;">No variants yet</li>';
            return;
        }
        
        listEl.innerHTML = variants.map(name => `
            <li>
                <span class="variant-name">📎 ${name}</span>
                <div class="variant-actions">
                    <button onclick="viewVariant('${name}')">View</button>
                    <button onclick="deleteVariant('${name}')">Delete</button>
                </div>
            </li>
        `).join('');
    } catch (error) {
        showError('Failed to load variants: ' + error.message);
    }
}

async function createVariant() {
    if (!currentPath || currentIsDir) return;
    
    const name = document.getElementById('variantName').value.trim();
    const content = document.getElementById('variantContent').value;
    
    if (!name) {
        showError('Please enter variant name');
        return;
    }
    
    try {
        await fetch(`${API_BASE}/variants${currentPath}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, content })
        });
        
        document.getElementById('variantName').value = '';
        document.getElementById('variantContent').value = '';
        showSuccess('Variant created successfully');
        loadVariants();
    } catch (error) {
        showError('Failed to create variant: ' + error.message);
    }
}

async function viewVariant(name) {
    try {
        const data = await fetch(`${API_BASE}/file${currentPath}/vo/${name}`).then(r => r.json());
        alert(`Variant: ${name}\n\nContent:\n${data.content}`);
    } catch (error) {
        showError('Failed to view variant: ' + error.message);
    }
}

async function deleteVariant(name) {
    if (!confirm(`Delete variant "${name}"?`)) return;
    
    try {
        await fetch(`${API_BASE}/variants${currentPath}?name=${encodeURIComponent(name)}`, {
            method: 'DELETE'
        });
        
        showSuccess('Variant deleted');
        loadVariants();
    } catch (error) {
        showError('Failed to delete variant: ' + error.message);
    }
}

// Search
async function performSearch() {
    const pattern = document.getElementById('searchPattern').value;
    const path = document.getElementById('searchPath').value || '/';
    const tagsStr = document.getElementById('searchTags').value;
    const description = document.getElementById('searchDescription').value;
    const recursive = document.getElementById('searchRecursive').checked;
    
    const params = new URLSearchParams();
    if (pattern) params.append('pattern', pattern);
    params.append('path', path);
    params.append('recursive', recursive);
    if (tagsStr) params.append('tags', tagsStr);
    if (description) params.append('description', description);
    
    try {
        const results = await fetch(`${API_BASE}/search?${params}`).then(r => r.json());
        
        const resultsEl = document.getElementById('searchResults');
        
        if (!results.paths || results.paths.length === 0) {
            resultsEl.innerHTML = '<div style="color: #999;">No results found</div>';
            return;
        }
        
        resultsEl.innerHTML = `
            <div class="result-count">${results.totalCount} result(s) found</div>
            ${results.paths.map(path => `
                <div class="result-item" onclick="selectFile('${path}', false)">
                    📄 ${path}
                </div>
            `).join('')}
        `;
    } catch (error) {
        showError('Search failed: ' + error.message);
    }
}

// Search files with improved error handling and UI
async function searchFiles() {
    console.log('[Search] Starting search...');
    
    const pattern = document.getElementById('searchPattern')?.value || '';
    const path = document.getElementById('searchPath')?.value || '/';
    const recursive = document.getElementById('searchRecursive')?.checked || true;
    const tags = document.getElementById('searchTags')?.value || '';
    const description = document.getElementById('searchDescription')?.value || '';
    
    console.log('[Search] Parameters:', { pattern, path, recursive, tags, description });
    
    // Validate at least one search criteria
    if (!pattern && !tags && !description) {
        showError('Please enter at least one search criteria (pattern, tags, or description)');
        return;
    }
    
    const resultsDiv = document.getElementById('searchResults');
    resultsDiv.innerHTML = '<p>🔄 Searching...</p>';
    
    try {
        // Build query parameters
        const params = new URLSearchParams();
        if (pattern) params.append('pattern', pattern);
        params.append('path', path);
        params.append('recursive', recursive.toString());
        if (tags) params.append('tags', tags);
        if (description) params.append('description', description);
        
        const url = `${API_BASE}/search?${params}`;
        console.log('[Search] Fetching:', url);
        
        const response = await fetch(url);
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP ${response.status}: ${errorText}`);
        }
        
        const data = await response.json();
        console.log('[Search] Response:', data);
        
        // Handle different response formats
        const results = data.Paths || data.paths || data || [];
        const totalCount = data.TotalCount || data.totalCount || results.length;
        
        console.log('[Search] Parsed results:', results, 'Total:', totalCount);
        
        // Display results
        if (!results || results.length === 0) {
            resultsDiv.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #999;">
                    <div style="font-size: 3em;">🔍</div>
                    <p>No files found matching your criteria</p>
                    <p style="font-size: 0.9em;">Try different search terms or enable recursive search</p>
                </div>
            `;
        } else {
            resultsDiv.innerHTML = `
                <div style="margin-bottom: 15px; padding: 10px; background: #e8f5e9; border-radius: 5px; border-left: 4px solid #4caf50;">
                    <strong style="color: #2e7d32;">✓ Found ${results.length} file(s)</strong>
                </div>
                <div style="max-height: 400px; overflow-y: auto;">
                    ${results.map(filePath => {
                        // Normalize path - remove double slashes
                        const normalizedPath = filePath.replace(/\/+/g, '/');
                        return `
                        <div class="search-result-item" 
                             onclick="selectFile('${normalizedPath}', false); showTab('editor')"
                             style="cursor: pointer; padding: 12px; margin: 8px 0; background: #f9f9f9; border-radius: 5px; border-left: 3px solid #667eea; transition: all 0.2s;"
                             onmouseover="this.style.background='#e8eaf6'; this.style.transform='translateX(5px)'"
                             onmouseout="this.style.background='#f9f9f9'; this.style.transform='translateX(0)'">
                            <div style="display: flex; align-items: center; gap: 10px;">
                                <span style="font-size: 1.2em;">📄</span>
                                <span style="font-family: monospace; color: #333;">${normalizedPath}</span>
                            </div>
                        </div>
                        `;
                    }).join('')}
                </div>
            `;
            
            showSuccess(`Found ${results.length} file(s)`);
        }
    } catch (error) {
        console.error('[Search] Error:', error);
        resultsDiv.innerHTML = `
            <div style="text-align: center; padding: 40px; color: #d32f2f;">
                <div style="font-size: 3em;">⚠️</div>
                <p><strong>Search Failed</strong></p>
                <p style="font-size: 0.9em; color: #666;">${error.message}</p>
            </div>
        `;
        showError('Search failed: ' + error.message);
    }
}

// Events
function toggleEvents() {
    if (isListeningToEvents) {
        stopEvents();
    } else {
        startEvents();
    }
}

function startEvents() {
    if (eventsSource) return;
    
    eventsSource = new EventSource(`${API_BASE}/events`);
    isListeningToEvents = true;
    document.getElementById('eventsBtn').textContent = '⏸️ Stop Listening';
    
    eventsSource.onmessage = (e) => {
        const event = JSON.parse(e.data);
        logEvent(event);
    };
    
    eventsSource.onerror = () => {
        showError('Event stream error');
        stopEvents();
    };
}

function stopEvents() {
    if (eventsSource) {
        eventsSource.close();
        eventsSource = null;
    }
    isListeningToEvents = false;
    document.getElementById('eventsBtn').textContent = '▶️ Start Listening';
}

function logEvent(event) {
    const log = document.getElementById('eventsLog');
    const types = ['Create', 'Modify', 'Delete', 'Rename', 'Chmod', 'Chown', 'Metadata', 'Variant'];
    const type = types[event.Type] || 'Unknown';
    const time = new Date().toLocaleTimeString();
    
    const entry = document.createElement('div');
    entry.className = 'event-entry';
    entry.innerHTML = `
        <span class="event-time">[${time}]</span>
        <span class="event-type event-${type.toLowerCase()}">${type}</span>
        <span class="event-path">${event.Path}</span>
    `;
    
    log.insertBefore(entry, log.firstChild);
    
    // Keep only last 100 events
    while (log.children.length > 100) {
        log.removeChild(log.lastChild);
    }
}

function clearEvents() {
    document.getElementById('eventsLog').innerHTML = '';
}

// Utility functions
function formatSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

function showSuccess(message) {
    showStatus(message, false);
}

function showError(message) {
    showStatus(message, true);
    console.error(message);
}

function showStatus(message, isError) {
    const el = document.createElement('div');
    el.className = `status-message ${isError ? 'error' : ''}`;
    el.textContent = message;
    document.body.appendChild(el);
    
    setTimeout(() => {
        el.style.opacity = '0';
        setTimeout(() => el.remove(), 300);
    }, 3000);
}
