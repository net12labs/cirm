// filepath: /home/lxk/Desktop/cirm/bins/vf_test/main.go
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lxk/cirm/bins/vfsql"
)

func main() {
	fmt.Println("=== VFSQL Test Application ===\n")

	// Clean up any existing database
	os.Remove("test.db")

	// Create database
	db, err := vfsql.Create("test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("✓ Created database: test.db")

	// Create a volume
	vol, err := db.CreateVolume("testvolume")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created volume: testvolume")

	// Get the default VFS
	vfs, err := vol.GetVFS("default")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Got default VFS\n")

	// Test 1: Basic File Operations
	fmt.Println("--- Test 1: Basic File Operations ---")
	testBasicFileOps(vfs)

	// Test 2: Directory Operations
	fmt.Println("\n--- Test 2: Directory Operations ---")
	testDirectoryOps(vfs)

	// Test 3: Metadata (Tags & Descriptions)
	fmt.Println("\n--- Test 3: Metadata Operations ---")
	testMetadata(vfs)

	// Test 4: File Variants
	fmt.Println("\n--- Test 4: File Variants ---")
	testVariants(vfs)

	// Test 5: Search Operations
	fmt.Println("\n--- Test 5: Search Operations ---")
	testSearch(vfs)

	// Test 6: Event Subscription
	fmt.Println("\n--- Test 6: Event Subscription ---")
	testEvents(vfs)

	fmt.Println("\n=== All Tests Completed Successfully! ===")
	fmt.Println("\nDatabase saved as: test.db")
	fmt.Println("You can inspect it with: sqlite3 test.db")
}

func testBasicFileOps(vfs *vfsql.VFS) {
	// Create a file
	f, err := vfs.Create("/hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Write to it
	content := "Hello, VFSQL! This is a test file.\n"
	n, err := f.Write([]byte(content))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Wrote %d bytes to /hello.txt\n", n)
	f.Close()

	// Read it back
	f, err = vfs.Open("/hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	fmt.Printf("  Read back: %q\n", string(data))

	// Get file info
	info, err := vfs.Stat("/hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  File info: size=%d, mode=%v\n", info.Size(), info.Mode())
}

func testDirectoryOps(vfs *vfsql.VFS) {
	// Create nested directories
	err := vfs.MkdirAll("/projects/2024/reports", 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Created directory structure: /projects/2024/reports")

	// Create files in different directories
	files := []string{
		"/projects/2024/report1.txt",
		"/projects/2024/report2.txt",
		"/projects/2024/reports/q1.txt",
		"/projects/2024/reports/q2.txt",
	}

	for _, path := range files {
		f, _ := vfs.Create(path)
		f.Write([]byte("Sample content for " + path))
		f.Close()
	}
	fmt.Printf("  Created %d test files\n", len(files))

	// List directory contents
	entries, err := vfs.ReadDir("/projects/2024")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Contents of /projects/2024: %d items\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("    - %s (dir: %v)\n", entry.Name(), entry.IsDir())
	}
}

func testMetadata(vfs *vfsql.VFS) {
	// Set description
	err := vfs.SetDescription("/hello.txt", "A friendly greeting file")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Set description on /hello.txt")

	// Set tags
	err = vfs.SetTags("/hello.txt", []string{"greeting", "test", "important"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Set tags: greeting, test, important")

	// Add another tag
	err = vfs.AddTag("/hello.txt", "reviewed")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Added tag: reviewed")

	// Get tags back
	tags, err := vfs.GetTags("/hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Current tags: %v\n", tags)

	// Get description
	desc, err := vfs.GetDescription("/hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Description: %q\n", desc)

	// Tag more files
	vfs.SetTags("/projects/2024/report1.txt", []string{"2024", "report", "draft"})
	vfs.SetTags("/projects/2024/report2.txt", []string{"2024", "report", "final"})
	vfs.SetTags("/projects/2024/reports/q1.txt", []string{"2024", "quarterly", "q1"})
	fmt.Println("  Tagged additional files")
}

func testVariants(vfs *vfsql.VFS) {
	// Create a mock image file
	imageData := []byte("ORIGINAL_IMAGE_DATA_1920x1080")
	f, _ := vfs.Create("/photo.jpg")
	f.Write(imageData)
	f.Close()
	fmt.Println("  Created /photo.jpg")

	// Create variants
	err := vfs.CreateVariant("/photo.jpg", "thumbnail.jpg", []byte("THUMB_200x200"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Created variant: thumbnail.jpg")

	err = vfs.CreateVariant("/photo.jpg", "medium.jpg", []byte("MEDIUM_800x600"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Created variant: medium.jpg")

	// List variants
	variants, err := vfs.ListVariants("/photo.jpg")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Variants: %v\n", variants)

	// Access variant via virtual path
	f, err = vfs.Open("/photo.jpg/vo/thumbnail.jpg")
	if err != nil {
		log.Fatal(err)
	}
	data, _ := io.ReadAll(f)
	f.Close()
	fmt.Printf("  Read thumbnail variant: %q\n", string(data))
}

func testSearch(vfs *vfsql.VFS) {
	// Simple name search
	results, err := vfs.FindByName("*.txt", &vfsql.FindOptions{
		BasePath:  "/",
		Recursive: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d .txt files\n", len(results))

	// Tag search
	results, err = vfs.FindByTag([]string{"2024"}, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d files with '2024' tag\n", len(results))

	// Advanced search
	searchResults, err := vfs.Search(&vfsql.SearchQuery{
		BasePath:    "/projects",
		Recursive:   true,
		Type:        vfsql.FileTypeFile,
		Tags:        []string{"2024", "report"},
		TagMatchAll: true,
		SortBy:      vfsql.SortByName,
		SortOrder:   vfsql.Ascending,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Advanced search found %d files:\n", len(searchResults.Paths))
	for _, path := range searchResults.Paths {
		fmt.Printf("    - %s\n", path)
	}
}

func testEvents(vfs *vfsql.VFS) {
	// Subscribe to events
	sub, err := vfs.Subscribe(&vfsql.EventFilter{
		Paths:     []string{"/"},
		Recursive: true,
		EventTypes: []vfsql.EventType{
			vfsql.EventCreate,
			vfsql.EventModify,
			vfsql.EventDelete,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer vfs.Unsubscribe(sub)

	// Event counter
	eventCount := 0
	done := make(chan bool)

	// Process events in background
	go func() {
		timeout := time.After(500 * time.Millisecond)
		for {
			select {
			case event := <-sub.Events:
				eventCount++
				eventType := []string{"Create", "Modify", "Delete", "Rename", "Chmod", "Chown", "Metadata", "Variant"}[event.Type]
				fmt.Printf("  Event #%d: %s - %s\n", eventCount, eventType, event.Path)
			case err := <-sub.Errors:
				fmt.Printf("  Error: %v\n", err)
			case <-timeout:
				done <- true
				return
			}
		}
	}()

	// Trigger some events
	time.Sleep(50 * time.Millisecond)

	// Create
	f, _ := vfs.Create("/event-test.txt")
	f.Write([]byte("Testing events"))
	f.Close()

	// Modify
	f, _ = vfs.OpenFile("/event-test.txt", os.O_RDWR, 0)
	f.Write([]byte(" - modified!"))
	f.Close()

	// Metadata
	vfs.SetTags("/event-test.txt", []string{"event-test"})

	// Delete
	vfs.Remove("/event-test.txt")

	// Wait for events to be processed
	<-done
	fmt.Printf("  Total events captured: %d\n", eventCount)
}
