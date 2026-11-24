package tickerapp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/google/uuid"
)

type FileCollector struct {
	Files []models.File
}

// fileFromString creates a File struct from a given path and adds it to the Files slice
func (fc *FileCollector) fileFromString(path string, info os.FileInfo) models.File {
	baseName := filepath.Base(path)

	file := models.File{
		Uuid:         uuid.NewString(),
		Name:         baseName,
		Path:         path,
		Size:         info.Size(),
		LastModified: info.ModTime(),
	}

	fc.Files = append(fc.Files, file)
	return file
}

func (fc *FileCollector) NumberOfFilesInFolder() int {
	return len(fc.Files)
}

// CollectFilesInFolder retrieves all the log files in the root directory
// and returns them as a slice of File structs
func (fc *FileCollector) CollectFilesInFolder(rootDir, path string) ([]models.File, error) {
	if path == "" || path == "/" {
		path = "data"
	}

	files := make([]models.File, 0, 16)

	trimmedPath := strings.TrimSuffix(path, "/")

	if trimmedPath == "" {
		trimmedPath = "data"
	}

	fullpath, err := filepath.Abs(filepath.Join(rootDir, trimmedPath))
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(fullpath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".log" {
			return nil
		}

		if fileInfo, err := d.Info(); err != nil {
			return err
		} else {
			files = append(files, fc.fileFromString(path, fileInfo))
			return nil
		}
	})

	return files, err
}
