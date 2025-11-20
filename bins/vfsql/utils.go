// filepath: /home/lxk/Desktop/cirm/bins/vfsql/utils.go
package vfsql

import (
	"path/filepath"
	"strings"
	"time"
)

// currentTimestamp returns the current Unix timestamp
func currentTimestamp() int64 {
	return time.Now().Unix()
}

// normalizePath cleans and normalizes a path
func normalizePath(path string) string {
	// Clean the path
	path = filepath.Clean(path)

	// Ensure it starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

// splitPath splits a path into directory and base name
func splitPath(path string) (string, string) {
	path = normalizePath(path)

	if path == "/" {
		return "", "/"
	}

	dir, base := filepath.Split(path)
	dir = normalizePath(dir)

	return dir, base
}

// joinPath joins path components
func joinPath(parts ...string) string {
	return normalizePath(filepath.Join(parts...))
}

// isVariantPath checks if a path contains the /vo/ virtual directory
func isVariantPath(path string) bool {
	return strings.Contains(path, "/vo/")
}

// splitVariantPath splits a variant path into original path and variant name
// e.g., "/images/photo.jpg/vo/thumbnail.jpg" -> ("/images/photo.jpg", "thumbnail.jpg")
func splitVariantPath(path string) (string, string, bool) {
	parts := strings.Split(path, "/vo/")
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// parseTagsString converts comma-separated tags string to slice
func parseTagsString(tags string) []string {
	if tags == "" {
		return nil
	}

	parts := strings.Split(tags, ",")
	result := make([]string, 0, len(parts))
	for _, tag := range parts {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// formatTagsString converts slice of tags to comma-separated string
func formatTagsString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	unique := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			unique = append(unique, tag)
		}
	}

	return strings.Join(unique, ",")
}

// hasTag checks if a tag is present in a tags string
func hasTag(tagsStr, tag string) bool {
	tags := parseTagsString(tagsStr)
	tag = strings.TrimSpace(tag)
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// addTag adds a tag to a tags string if not already present
func addTag(tagsStr, tag string) string {
	tags := parseTagsString(tagsStr)
	tag = strings.TrimSpace(tag)

	// Check if already exists
	for _, t := range tags {
		if t == tag {
			return tagsStr
		}
	}

	tags = append(tags, tag)
	return formatTagsString(tags)
}

// removeTag removes a tag from a tags string
func removeTag(tagsStr, tag string) string {
	tags := parseTagsString(tagsStr)
	tag = strings.TrimSpace(tag)

	result := make([]string, 0, len(tags))
	for _, t := range tags {
		if t != tag {
			result = append(result, t)
		}
	}

	return formatTagsString(result)
}

// globToLike converts a glob pattern to SQL LIKE pattern
func globToLike(pattern string) string {
	// Replace * with % and ? with _
	pattern = strings.ReplaceAll(pattern, "*", "%")
	pattern = strings.ReplaceAll(pattern, "?", "_")
	return pattern
}
