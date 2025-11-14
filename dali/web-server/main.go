package webserver

import (
	"fmt"
	"net/http"
)

type WebServer struct {
	// WebServer fields here
	Port      int
	OnRequest func(req *Request) error
	mux       *http.ServeMux
}

func (ws *WebServer) Start() error {
	// Start the web server
	if ws.Port == 0 {
		ws.Port = 8080
	}

	if ws.mux == nil {
		ws.mux = http.NewServeMux()
	}

	addr := fmt.Sprintf(":%d", ws.Port)
	fmt.Printf("Starting web server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, ws.mux)
}

func (ws *WebServer) AddRoute(path string, handler func(req *Request)) {
	// Add a new route to the web server
	if ws.mux == nil {
		ws.mux = http.NewServeMux()
	}

	fmt.Printf("Registering route: %s\n", path)
	ws.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(w, r)
		// fmt.Printf("Handler triggered for path: %s (requested: %s)\n", path, r.URL.Path)
		handler(req)
	})
}

func NewWebServer() *WebServer {
	return &WebServer{
		mux: http.NewServeMux(),
	}
}

func (ws *WebServer) Init() error {
	// Initialize the web server
	return nil
}
