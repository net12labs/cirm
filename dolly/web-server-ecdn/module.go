package ecdnwebserver

import (
	"embed"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	webserver "github.com/net12labs/cirm/mali/web-server"
)

type Server struct {
	WebServer   *webserver.WebServer
	Fs          *embed.FS
	UrlBasePath string
	FsBasePath  string
}

func NewServer() *Server {
	return &Server{}
}

// this should be serving js, fonts, css, images, static json etc.
// properly set up mime types
// extract the path after Url Base path and append to FsBasePath to get file path in Fs

func (s *Server) Init() error {
	if s.WebServer == nil {
		return fmt.Errorf("WebServer is not set")
	}

	// Ensure UrlBasePath ends without trailing slash for consistent path handling
	urlBase := strings.TrimSuffix(s.UrlBasePath, "/")
	fsBase := strings.TrimSuffix(s.FsBasePath, "/")

	println("ECDN Server Init - UrlBasePath:", s.UrlBasePath, "FsBasePath:", s.FsBasePath)

	// Add wildcard route to handle all paths under the base path
	// Using Go 1.22+ ServeMux wildcard pattern
	s.WebServer.AddRoute(urlBase+"/", func(req *webserver.Request) {

		// Extract the path after the base URL
		requestPath := req.Path.Path
		relativePath := strings.TrimPrefix(requestPath, urlBase)
		relativePath = strings.TrimPrefix(relativePath, "/")

		// Construct the full file path in the embedded FS
		filePath := fsBase + "/" + relativePath

		// Handle default index.html for directory paths
		if strings.HasSuffix(filePath, "/") || relativePath == "" {
			filePath = fsBase + "/index.html"
		}

		// Read the file from embedded FS
		data, err := s.Fs.ReadFile(filePath)
		if err != nil {
			println("ECDN file not found:", filePath, "Error:", err.Error())
			req.WriteResponse404()
			return
		}

		// Determine MIME type based on file extension
		ext := filepath.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			// Default to octet-stream if MIME type cannot be determined
			mimeType = "application/octet-stream"
		}

		// Set Content-Type header and write response
		req.Response.Headers.Set("Content-Type", mimeType)
		req.WriteResponse(data)
	})

	return nil
}
