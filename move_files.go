package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)


// Define the file type to directory mapping
var fileTypeMapping = map[string]string{
	".jpg":  "images",
	".jpeg": "images",
	".png":  "images",
	".gif":  "images",
	".txt":  "Documents",
	".pdf":  "pdfs",
	".docx": "Documents",
	".xlsx": "Documents",
	".mp4":  "Videos",
	".mov":  "Videos",
	".mp3":  "music",
	".wav":  "music",
	".m4a":  "music",
	".csv":  "csv",
	".md":  "markdown",
	// Add more mappings as needed
}

func moveFilesToNAS(desktopPath string, nasBasePath string, fileTypeMapping map[string]string) error {
	// Ensure the NAS base path exists
	if _, err := os.Stat(nasBasePath); os.IsNotExist(err) {
		return fmt.Errorf("NAS path '%s' does not exist", nasBasePath)
	}

	// Open the Desktop directory
	dir, err := os.Open(desktopPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get the list of files in the Desktop directory
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Get the file path
		filePath := filepath.Join(desktopPath, file.Name())

		// Get the file extension
		fileExtension := strings.ToLower(filepath.Ext(file.Name()))

		// Determine the target directory based on the file type
		targetDir, exists := fileTypeMapping[fileExtension]
		if !exists {
			fmt.Printf("No mapping for file type: %s for file %s. Skipping.\n", fileExtension, file.Name())
			continue
		}

		targetPath := filepath.Join(nasBasePath, targetDir)
		
		// Ensure the target directory exists on the NAS
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return err
		}

		// Copy the file to the target directory
		destination := filepath.Join(targetPath, file.Name())
		if err := copyFile(filePath, destination); err != nil {
			return err
		}

		// Delete the original file
		if err := os.Remove(filePath); err != nil {
			return err
		}

		fmt.Printf("Moved %s to %s\n", file.Name(), targetPath)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return destinationFile.Sync()
}

func main() {
	desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop")
	nasBasePath := "/Volumes/Pioneer/archive"

	if err := moveFilesToNAS(desktopPath, nasBasePath, fileTypeMapping); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
