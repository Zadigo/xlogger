package test

import (
	"path/filepath"
	"testing"

	"github.com/Zadigo/xlogger_backend/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestFilePathWalk(t *testing.T) {
	path, err := filepath.Abs("../data")
	assert.NoError(t, err, "Should not error when getting absolute path")

	t.Run("List log files", func(t *testing.T) {
		logFiles, err := backend.FilePathWalkDir(path)
		assert.NoError(t, err, "Should not error when walking the directory")
		assert.Len(t, logFiles, 3, "Should find 3 log files")
	})
}
