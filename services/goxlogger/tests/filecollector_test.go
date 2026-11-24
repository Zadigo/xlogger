package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/stretchr/testify/assert"
)

func TestFileCollector(t *testing.T) {
	t.Run("should collect files in folder", func(t *testing.T) {
		t.Parallel()

		fc := &tickerapp.FileCollector{}
		files, err := fc.CollectFilesInFolder("../", "data")
		assert.NoError(t, err)
		assert.Greater(t, len(files), 0)
	})

	t.Run("should collect files in folder with empty path", func(t *testing.T) {
		t.Parallel()

		fc := &tickerapp.FileCollector{}
		files, err := fc.CollectFilesInFolder("../", "")
		assert.NoError(t, err)
		assert.Greater(t, len(files), 0)
	})
}

func BenchmarkTestFileo(b *testing.B) {
	b.ResetTimer()
	b.SetParallelism(30)

	for b.Loop() {
		fc := &tickerapp.FileCollector{}
		_, err := fc.CollectFilesInFolder("../", "data")
		if err != nil {
			b.Fatal(err)
		}
	}
}
