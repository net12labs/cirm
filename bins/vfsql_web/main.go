// filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lxk/cirm/bins/vfsql/api"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "vfsql.db", "Path to database file")
	port := flag.Int("port", 8080, "HTTP server port")
	webDir := flag.String("web", "../vfsql/web", "Path to web directory")
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║      VFSQL Web Interface Server         ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Database: %s\n", *dbPath)
	fmt.Printf("Port:     %d\n", *port)
	fmt.Printf("Web Dir:  %s\n", *webDir)
	fmt.Println()

	// Create symlink to web directory if it doesn't exist
	if _, err := os.Stat("web"); os.IsNotExist(err) {
		if _, err := os.Stat(*webDir); err == nil {
			fmt.Println("Creating symlink to web directory...")
			os.Symlink(*webDir, "web")
		}
	}

	// Create server
	server, err := api.NewServer(*dbPath, *port)
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}

	fmt.Println("Web Interface: http://localhost:8080")
	fmt.Println("API Endpoint:  http://localhost:8080/api")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println()

	// Start server
	if err := server.Start(); err != nil {
		log.Fatal("Server error: ", err)
	}
}
