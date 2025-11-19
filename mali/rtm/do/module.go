package do

import (
	"os"
	"path/filepath"
)

type Do struct {
	// Do fields here
}

func NewDo() *Do {
	return &Do{}
}

func InitFsPath(path string) error {
	// Initialize file system path
	// make directory path recursively. Split off the last part if it is a file
	dir := path

	// Check if path looks like a file (has an extension)
	if filepath.Ext(path) != "" {
		dir = filepath.Dir(path)
	}

	// Create directory recursively
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return nil
}
