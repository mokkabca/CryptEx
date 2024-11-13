package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Decode the .bot file and restore it with the original extension
func decodeFile(encodedFilePath string) error {
	// Read the encoded .bot file
	encodedData, err := ioutil.ReadFile(encodedFilePath)
	if err != nil {
		return fmt.Errorf("error reading encoded file %s: %v", encodedFilePath, err)
	}

	// Split the first line (extension) and the rest (base64 content)
	lines := strings.SplitN(string(encodedData), "\n", 2)
	if len(lines) < 2 {
		return fmt.Errorf("invalid .bot file format: missing extension or base64 content")
	}

	originalExtension := lines[0]
	encodedContent := lines[1]

	// Check if the extension is valid (e.g., .jpg, .txt, .png, etc.)
	if !strings.HasPrefix(originalExtension, ".") {
		return fmt.Errorf("invalid file extension in .bot file: %s", originalExtension)
	}

	// Decode the base64 content
	decodedData, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		return fmt.Errorf("error decoding base64 data: %v", err)
	}

	// Restore the original file with its extension
	originalFileName := strings.TrimSuffix(filepath.Base(encodedFilePath), ".bot") + originalExtension
	err = ioutil.WriteFile(originalFileName, decodedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing decoded file %s: %v", originalFileName, err)
	}

	log.Printf("Successfully restored original file: %s", originalFileName)

	// Delete the original .bot file after decoding
	err = os.Remove(encodedFilePath)
	if err != nil {
		return fmt.Errorf("error deleting .bot file %s: %v", encodedFilePath, err)
	}

	log.Printf("Successfully deleted .bot file: %s", encodedFilePath)
	return nil
}

func main() {
	// Replace with your directory containing the .bot files
	dir := "./" // Or specify the directory where your .bot files are located

	// List all .bot files in the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Decode each .bot file
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bot" {
			encodedFilePath := filepath.Join(dir, file.Name())

			// Decode the .bot file and restore it
			err := decodeFile(encodedFilePath)
			if err != nil {
				log.Printf("Error decoding file %s: %v", file.Name(), err)
			}
		}
	}
}
