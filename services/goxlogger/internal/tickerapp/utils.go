package tickerapp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


type FileCollector struct {
	Files []File
}

// fileFromString creates a File struct from a given path and adds it to the Files slice
func (fc *FileCollector) fileFromString(path string) File {
	baseName := filepath.Base(path)

	file := File{Name: baseName, Path: path}
	fc.Files = append(fc.Files, file)

	return file
}

func (fc *FileCollector) NumberOfFilesInFolder() int {
	return len(fc.Files)
}

// CollectFilesInFolder retrieves all the log files in the root directory
// and returns them as a slice of File structs
func (fc *FileCollector) CollectFilesInFolder(rootDir, path string) ([]File, error) {
	var files []File

	trimmedPath := strings.TrimSuffix(path, "/")

	if trimmedPath == "" {
		trimmedPath = "data"
	}

	fullpath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootDir, trimmedPath))
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".log" {
			// log.Printf("⚠️ Skipping file %s:", path)
			return nil
		}

		if !info.IsDir() {
			files = append(files, fc.fileFromString(path))
		}

		return nil
	})
	return files, err
}
