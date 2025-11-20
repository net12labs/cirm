// filepath: /home/lxk/Desktop/cirm/bins/vfsql/api/server.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lxk/cirm/bins/vfsql"
)

// Server represents the HTTP API server
type Server struct {
	db   *vfsql.DB
	vfs  *vfsql.VFS
	mux  *http.ServeMux
	port int
}

// NewServer creates a new API server
func NewServer(dbPath string, port int) (*Server, error) {
	db, err := vfsql.Open(dbPath)
	if err != nil {
		return nil, err
	}

	// Get or create default volume
	vol, err := db.GetVolume("default")
	if err != nil {
		vol, err = db.CreateVolume("default")
		if err != nil {
			return nil, err
		}
	}

	// Get default VFS
	vfs, err := vol.GetVFS("default")
	if err != nil {
		return nil, err
	}

	s := &Server{
		db:   db,
		vfs:  vfs,
		mux:  http.NewServeMux(),
		port: port,
	}

	s.setupRoutes()
	return s, nil
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// API routes first (more specific)
	s.mux.HandleFunc("/api/upload", s.handleUpload)
	s.mux.HandleFunc("/api/download/", s.handleDownload)
	s.mux.HandleFunc("/api/test-create", s.handleTestCreate)
	s.mux.HandleFunc("/api/files", s.handleFiles)
	s.mux.HandleFunc("/api/file/", s.handleFile)
	s.mux.HandleFunc("/api/dirs", s.handleDirs)
	s.mux.HandleFunc("/api/dir/", s.handleDir)
	s.mux.HandleFunc("/api/metadata/", s.handleMetadata)
	s.mux.HandleFunc("/api/tags/", s.handleTags)
	s.mux.HandleFunc("/api/variants/", s.handleVariants)
	s.mux.HandleFunc("/api/search", s.handleSearch)
	s.mux.HandleFunc("/api/stat/", s.handleStat)
	s.mux.HandleFunc("/api/tree", s.handleTree)
	s.mux.HandleFunc("/api/events", s.handleEvents)

	// Serve static files (must be last)
	s.mux.HandleFunc("/", s.handleStatic)
}

// handleStatic serves the web interface
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// If requesting root, serve index.html
	if r.URL.Path == "/" {
		s.serveFile(w, r, "index.html")
		return
	}

	// Serve other static files
	s.serveFile(w, r, r.URL.Path[1:])
}

// serveFile serves static files from embedded or filesystem
func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, filename string) {
	// Try to find web files in different locations
	paths := []string{
		"web/" + filename,
		"../vfsql/web/" + filename,
		"../../vfsql/web/" + filename,
	}

	var content []byte
	var err error

	for _, path := range paths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		// If file not found, return 404
		if filename != "index.html" {
			http.NotFound(w, r)
			return
		}

		// Serve embedded fallback HTML for index
		s.serveFallbackHTML(w)
		return
	}

	// Set content type based on file extension
	contentType := "text/plain"
	if strings.HasSuffix(filename, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(filename, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "application/javascript"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

// serveFallbackHTML serves a minimal HTML page if web files aren't found
func (s *Server) serveFallbackHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html>
<head>
    <title>VFSQL API Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #667eea; }
        .endpoint { background: #f9f9f9; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { display: inline-block; padding: 5px 10px; border-radius: 3px; font-weight: bold; margin-right: 10px; }
        .get { background: #4caf50; color: white; }
        .post { background: #2196f3; color: white; }
        .put { background: #ff9800; color: white; }
        .delete { background: #f44336; color: white; }
        code { background: #eee; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🗂️ VFSQL API Server</h1>
        <p>The VFSQL API is running. Web files not found in expected locations.</p>
        
        <h2>API Endpoints</h2>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/files</code> - List files
        </div>
        
        <div class="endpoint">
            <span class="method post">POST</span>
            <code>/api/files</code> - Create file
            <br><small>Body: <code>{"path": "/test.txt", "content": "..."}</code></small>
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/file/:path</code> - Get file content
        </div>
        
        <div class="endpoint">
            <span class="method put">PUT</span>
            <code>/api/file/:path</code> - Update file
        </div>
        
        <div class="endpoint">
            <span class="method delete">DELETE</span>
            <code>/api/file/:path</code> - Delete file
        </div>
        
        <div class="endpoint">
            <span class="method post">POST</span>
            <code>/api/dirs</code> - Create directory
            <br><small>Body: <code>{"path": "/mydir"}</code></small>
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/dir/:path</code> - List directory
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/metadata/:path</code> - Get metadata
        </div>
        
        <div class="endpoint">
            <span class="method put">PUT</span>
            <code>/api/metadata/:path</code> - Set metadata
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/search</code> - Search files
            <br><small>Params: <code>?pattern=*.txt&recursive=true&tags=test</code></small>
        </div>
        
        <div class="endpoint">
            <span class="method get">GET</span>
            <code>/api/tree</code> - Get directory tree
        </div>
        
        <h2>Quick Test</h2>
        <p>Try creating a file:</p>
        <pre>curl -X POST http://localhost:8080/api/files \
  -H "Content-Type: application/json" \
  -d '{"path":"/test.txt","content":"Hello!"}'</pre>
        
        <p>Then read it:</p>
        <pre>curl http://localhost:8080/api/file/test.txt</pre>
        
        <h2>Setup Web Interface</h2>
        <p>To use the web interface, ensure the <code>web/</code> directory is accessible:</p>
        <ol>
            <li>Make sure web files exist at: <code>../vfsql/web/</code></li>
            <li>Or run the server from the correct directory</li>
            <li>Or copy web files to <code>./web/</code></li>
        </ol>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Starting VFSQL API server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.enableCORS(s.mux))
}

// enableCORS adds CORS headers
func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleFiles lists files in root or creates new file
func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		entries, err := s.vfs.ReadDir("/")
		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, entries)

	case "POST":
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		f, err := s.vfs.Create(req.Path)
		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		n, err := f.Write([]byte(req.Content))
		f.Close()

		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		fmt.Printf("Created file: path=%s, wrote=%d bytes\n", req.Path, n)

		s.jsonResponse(w, map[string]interface{}{
			"status": "created",
			"path":   req.Path,
			"bytes":  n,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleFile manages single file operations
func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/file")

	switch r.Method {
	case "GET":
		f, err := s.vfs.Open(path)
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]interface{}{
			"path":    path,
			"content": string(content),
		})

	case "PUT":
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		f, err := s.vfs.OpenFile(path, 2, 0644) // O_RDWR
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}
		defer f.Close()

		f.Truncate(0)
		if _, err := f.Write([]byte(req.Content)); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "updated"})

	case "DELETE":
		if err := s.vfs.Remove(path); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, map[string]string{"status": "deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDirs handles directory operations
func (s *Server) handleDirs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		if err := s.vfs.MkdirAll(req.Path, 0755); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "created", "path": req.Path})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDir handles single directory operations
func (s *Server) handleDir(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/dir")
	if path == "" {
		path = "/"
	}

	switch r.Method {
	case "GET":
		entries, err := s.vfs.ReadDir(path)
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}

		result := make([]map[string]interface{}, len(entries))
		for i, entry := range entries {
			info, _ := entry.Info()
			result[i] = map[string]interface{}{
				"name":  entry.Name(),
				"isDir": entry.IsDir(),
				"size":  info.Size(),
				"mode":  info.Mode().String(),
			}
		}

		s.jsonResponse(w, result)

	case "DELETE":
		if err := s.vfs.RemoveAll(path); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, map[string]string{"status": "deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleMetadata manages file metadata
func (s *Server) handleMetadata(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/metadata")

	switch r.Method {
	case "GET":
		desc, err := s.vfs.GetDescription(path)
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}

		tags, err := s.vfs.GetTags(path)
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}

		s.jsonResponse(w, map[string]interface{}{
			"description": desc,
			"tags":        tags,
		})

	case "PUT":
		var req struct {
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		if req.Description != "" {
			if err := s.vfs.SetDescription(path, req.Description); err != nil {
				s.jsonError(w, err, http.StatusInternalServerError)
				return
			}
		}

		if len(req.Tags) > 0 {
			if err := s.vfs.SetTags(path, req.Tags); err != nil {
				s.jsonError(w, err, http.StatusInternalServerError)
				return
			}
		}

		s.jsonResponse(w, map[string]string{"status": "updated"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTags manages tags
func (s *Server) handleTags(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/tags")

	switch r.Method {
	case "POST":
		var req struct {
			Tag string `json:"tag"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		if err := s.vfs.AddTag(path, req.Tag); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "added"})

	case "DELETE":
		tag := r.URL.Query().Get("tag")
		if tag == "" {
			s.jsonError(w, fmt.Errorf("tag parameter required"), http.StatusBadRequest)
			return
		}

		if err := s.vfs.RemoveTag(path, tag); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "removed"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleVariants manages file variants
func (s *Server) handleVariants(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/variants")

	switch r.Method {
	case "GET":
		variants, err := s.vfs.ListVariants(path)
		if err != nil {
			s.jsonError(w, err, http.StatusNotFound)
			return
		}
		s.jsonResponse(w, variants)

	case "POST":
		var req struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		if err := s.vfs.CreateVariant(path, req.Name, []byte(req.Content)); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "created"})

	case "DELETE":
		name := r.URL.Query().Get("name")
		if name == "" {
			s.jsonError(w, fmt.Errorf("name parameter required"), http.StatusBadRequest)
			return
		}

		if err := s.vfs.RemoveVariant(path, name); err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		s.jsonResponse(w, map[string]string{"status": "deleted"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSearch performs file search
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := &vfsql.SearchQuery{
		NamePattern: q.Get("pattern"),
		BasePath:    q.Get("path"),
		Recursive:   q.Get("recursive") == "true",
		Description: q.Get("description"),
	}

	if tags := q.Get("tags"); tags != "" {
		query.Tags = strings.Split(tags, ",")
	}

	if limit := q.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			query.Limit = l
		}
	}

	results, err := s.vfs.Search(query)
	if err != nil {
		s.jsonError(w, err, http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, results)
}

// handleStat returns file/directory stats
func (s *Server) handleStat(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/stat")

	info, err := s.vfs.Stat(path)
	if err != nil {
		s.jsonError(w, err, http.StatusNotFound)
		return
	}

	s.jsonResponse(w, map[string]interface{}{
		"name":    info.Name(),
		"size":    info.Size(),
		"mode":    info.Mode().String(),
		"modTime": info.ModTime().Format(time.RFC3339),
		"isDir":   info.IsDir(),
	})
}

// handleTree returns directory tree
func (s *Server) handleTree(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	tree, err := s.buildTree(path, 3) // max depth 3
	if err != nil {
		s.jsonError(w, err, http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, tree)
}

// buildTree recursively builds directory tree
func (s *Server) buildTree(path string, depth int) (map[string]interface{}, error) {
	if depth <= 0 {
		return nil, nil
	}

	info, err := s.vfs.Stat(path)
	if err != nil {
		return nil, err
	}

	node := map[string]interface{}{
		"name":  filepath.Base(path),
		"path":  path,
		"isDir": info.IsDir(),
		"size":  info.Size(),
	}

	if info.IsDir() {
		entries, err := s.vfs.ReadDir(path)
		if err != nil {
			return node, nil
		}

		children := []map[string]interface{}{}
		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			if child, err := s.buildTree(childPath, depth-1); err == nil && child != nil {
				children = append(children, child)
			}
		}
		node["children"] = children
	}

	return node, nil
}

// handleEvents handles SSE for real-time events
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sub, err := s.vfs.Subscribe(&vfsql.EventFilter{
		Paths:     []string{"/"},
		Recursive: true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer s.vfs.Unsubscribe(sub)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case event := <-sub.Events:
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// handleTestCreate creates a test file to verify everything works
func (s *Server) handleTestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testPath := "/test-file.txt"
	testContent := "This is a test file created by the API.\nIt has multiple lines.\nCurrent time: " + time.Now().Format(time.RFC3339)

	// Create file
	f, err := s.vfs.Create(testPath)
	if err != nil {
		s.jsonError(w, fmt.Errorf("failed to create: %w", err), http.StatusInternalServerError)
		return
	}

	// Write content
	n, err := f.Write([]byte(testContent))
	if err != nil {
		f.Close()
		s.jsonError(w, fmt.Errorf("failed to write: %w", err), http.StatusInternalServerError)
		return
	}

	f.Close()

	// Verify by reading back
	f2, err := s.vfs.Open(testPath)
	if err != nil {
		s.jsonError(w, fmt.Errorf("failed to open for verify: %w", err), http.StatusInternalServerError)
		return
	}

	readBack, err := io.ReadAll(f2)
	f2.Close()

	if err != nil {
		s.jsonError(w, fmt.Errorf("failed to read back: %w", err), http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, map[string]interface{}{
		"status":        "created",
		"path":          testPath,
		"written_bytes": n,
		"content_len":   len(testContent),
		"read_back_len": len(readBack),
		"content_match": string(readBack) == testContent,
		"read_back":     string(readBack),
	})
}

// handleDownload serves files for download with proper headers
func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/download")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get file info first
	info, err := s.vfs.Stat(path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}

	// Open file
	f, err := s.vfs.Open(path)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Read content
	content, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Debug: log the size
	fmt.Printf("Download: path=%s, size=%d, content_length=%d\n", path, len(content), info.Size())

	// Get filename
	filename := filepath.Base(path)

	// Determine content type based on extension
	contentType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		contentType = "text/plain; charset=utf-8"
	case ".json":
		contentType = "application/json"
	case ".xml":
		contentType = "application/xml"
	case ".html":
		contentType = "text/html; charset=utf-8"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".md":
		contentType = "text/markdown; charset=utf-8"
	case ".csv":
		contentType = "text/csv"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".pdf":
		contentType = "application/pdf"
	case ".zip":
		contentType = "application/zip"
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Header().Set("Cache-Control", "no-cache")

	// Write content
	written, err := w.Write(content)
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		return
	}
	fmt.Printf("Wrote %d bytes\n", written)
} // handleUpload handles multipart file uploads
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (32MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		s.jsonError(w, err, http.StatusBadRequest)
		return
	}

	targetPath := r.FormValue("path")
	if targetPath == "" {
		targetPath = "/"
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		s.jsonError(w, fmt.Errorf("no files provided"), http.StatusBadRequest)
		return
	}

	results := make([]map[string]interface{}, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			results = append(results, map[string]interface{}{
				"name":    fileHeader.Filename,
				"success": false,
				"error":   err.Error(),
			})
			continue
		}

		// Read file content
		content, err := io.ReadAll(file)
		file.Close()

		if err != nil {
			results = append(results, map[string]interface{}{
				"name":    fileHeader.Filename,
				"success": false,
				"error":   err.Error(),
			})
			continue
		}

		fmt.Printf("Upload: filename=%s, read=%d bytes\n", fileHeader.Filename, len(content))

		// Create file path
		filePath := filepath.Join(targetPath, fileHeader.Filename)
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}

		fmt.Printf("Upload: creating at path=%s\n", filePath)

		// Create file in VFS
		f, err := s.vfs.Create(filePath)
		if err != nil {
			fmt.Printf("Upload: create failed: %v\n", err)
			results = append(results, map[string]interface{}{
				"name":    fileHeader.Filename,
				"success": false,
				"error":   err.Error(),
			})
			continue
		}

		n, err := f.Write(content)
		f.Close()

		fmt.Printf("Upload: wrote %d bytes, err=%v\n", n, err)

		if err != nil {
			results = append(results, map[string]interface{}{
				"name":    fileHeader.Filename,
				"success": false,
				"error":   err.Error(),
			})
			continue
		}

		// Verify by reading back
		verifyF, verifyErr := s.vfs.Open(filePath)
		if verifyErr == nil {
			verifyContent, _ := io.ReadAll(verifyF)
			verifyF.Close()
			fmt.Printf("Upload: verification read back %d bytes\n", len(verifyContent))
		}

		results = append(results, map[string]interface{}{
			"name":    fileHeader.Filename,
			"success": true,
			"path":    filePath,
			"size":    len(content),
		})
	}

	s.jsonResponse(w, map[string]interface{}{
		"uploaded": len(files),
		"results":  results,
	})
}

// jsonResponse sends JSON response
func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// jsonError sends JSON error response
func (s *Server) jsonError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// Close closes the server
func (s *Server) Close() error {
	return s.db.Close()
}
