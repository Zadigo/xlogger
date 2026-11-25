package utils

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"path/filepath"
)

func GetAbsolutePath(strPath string) (string, error) {
	absPath, err := filepath.Abs(strPath)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get absolute path: %v", err)
	}

	result := path.Ext(absPath)
	if result != "" {
		return "", fmt.Errorf("❌ Base directory should be a directory, got a file: %s", absPath)
	}

	return absPath, nil
}

func GetRootDir(ctx context.Context) string {
	strRootDir := ctx.Value("rootDir")

	if strRootDir == nil {
		slog.Error("❌ Root directory is not set in the context")
		strRootDir = "."
	}

	rootDir, ok := strRootDir.(string)
	if !ok {
		slog.Error("❌ Root directory in context is not a string")
		strRootDir = "."
	}

	if rootDir == "" {
		slog.Error("❌ Root directory is empty")
		strRootDir = "."
	}

	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		slog.Error(fmt.Sprintf("❌ Failed to get absolute path of root directory: %s", err))

	}

	return absRootDir
}
