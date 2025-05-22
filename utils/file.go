package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// PATH specifies the base directory where files will be stored or accessed.
var PATH = "assets"

// UploadFile saves the uploaded file to the specified path, creating necessary directories if they don't exist.
// It takes a file header and a file path as parameters, returning an error if any operation fails.
func UploadFile(file *multipart.FileHeader, path string) error {
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		return fmt.Errorf("invalid path: %s", path)
	}
	fileID := parts[len(parts)-1]
	dirParts := parts[:len(parts)-1]
	dirPath := filepath.Join(PATH, filepath.Join(dirParts...))

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, fileID)

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer func(uploadedFile multipart.File) {
		err := uploadedFile.Close()
		if err != nil {
			panic(err)
		}
	}(uploadedFile)

	targetFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(targetFile *os.File) {
		err := targetFile.Close()
		if err != nil {
			panic(err)
		}
	}(targetFile)

	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		return err
	}

	return nil
}

// GetExtensions extracts and returns the file extension from a given filename after the last dot.
func GetExtensions(filename string) string {
	return strings.Split(filename, ".")[len(strings.Split(filename, "."))-1]
}
