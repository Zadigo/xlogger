package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/stretchr/testify/assert"
)

func TestFileCollector(t *testing.T) {
	fc := &tickerapp.FileCollector{}

	t.Run("should collect files in folder", func(t *testing.T) {
		files, err := fc.CollectFilesInFolder("../", "data")
		assert.NoError(t, err)
		assert.Greater(t, len(files), 0)
	})
}
