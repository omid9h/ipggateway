package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDirectory returns absolute path of caller + dir if directory is
// not present, it creates the directory
func GetDirectory(dir string) (string, error) {
	// Get the absolute path to the directory containing the source file
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to determine the current file path")
	}

	// Build the path to the "./dir/" directory
	targetDir := filepath.Join(currentDir, dir)

	// Check if the directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.Mkdir(targetDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create %s directory: %v", dir, err)
		}
	}

	return targetDir, nil
}
