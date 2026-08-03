package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Zadigo/goxlogger/internal/backend"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/stretchr/testify/assert"
)

func instanceFixture() *tickerapp.FileRedis {
	ctx := context.WithValue(context.Background(), "rootDir", "../")

	redisClient := backend.NewRedisBackend(ctx)
	filesRedis := tickerapp.NewFileRedis(ctx, redisClient)

	return filesRedis
}

func TestCollectFilesInFolder(t *testing.T) {
	filesRedis := instanceFixture()
	files, err := filesRedis.CollectFilesInFolder("/data")

	assert.Nil(t, err)
	assert.NotEmpty(t, files)
	assert.NoError(t, err)
	assert.True(t, len(files) > 0)

	for _, file := range files {
		fileInfo, err := os.Stat(file.Path)
		assert.Nil(t, err)
		assert.NotNil(t, fileInfo)
	}
}

func TestSaveFiles(t *testing.T) {
	t.Run("should save file", func(t *testing.T) {
		fileRedis := instanceFixture()

		testFile := []tickerapp.File{
			{
				Name: "example1.log",
				Path: filepath.Join("../data", "example1.log"),
			},
		}

		err := fileRedis.SaveFiles(testFile)
		assert.Nil(t, err)

		t.Cleanup(func() {
			redisClient := backend.NewRedisBackend(t.Context())
			redisClient.FlushAll(context.Background()).Err()
		})
	})
}

func TestReadFile(t *testing.T) {
	filesRedis := instanceFixture()

	path := t.TempDir()
	fullPath := filepath.Join(path, "example.log")
	os.WriteFile(fullPath, []byte(`146.70.194.73 - - [20/Oct/2025:20:50:30 +0000] "GET /blog/insights/les-9-meilleurs-generateurs-de-videos-par-ia HTTP/2.0" 200 38314 "-" "-" 1 "myprosite@docker" "http://172.21.0.18:8000" 89ms`), 0644)

	t.Run("should read file", func(t *testing.T) {
		logs, err := filesRedis.ReadFile(fullPath, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, logs)
		assert.True(t, len(logs) > 0)
	})
}

func TestCacheContent(t *testing.T) {
	fileRedis := instanceFixture()

	t.Run("should cache content", func(t *testing.T) {
		content := []string{"line1", "line2", "line3"}
		err := fileRedis.CacheLogs("example1.log", content)
		assert.Nil(t, err)

		cachedContent, err := fileRedis.GetLogs("example1.log")
		assert.Nil(t, err)
		assert.NotEmpty(t, cachedContent)
		// assert.Equal(t, []logic.LogLine{
		// 	{RawLine: "line1"},
		// 	{RawLine: "line2"},
		// 	{RawLine: "line3"},
		// }, cachedContent)
	})
}

func TestImplementation(t *testing.T) {
	ctx := context.WithValue(t.Context(), "rootDir", "../")

	redisClient := backend.NewRedisBackend(ctx)
	filesRedis := tickerapp.NewFileRedis(ctx, redisClient)

	t.Run("Should load files", func(t *testing.T) {
		files, err := filesRedis.CollectFilesInFolder("/data")
		assert.Nil(t, err)
		assert.NotEmpty(t, files)
		assert.NoError(t, err)
		assert.True(t, len(files) > 0)
	})

	t.Run("Should save files", func(t *testing.T) {
		files, _ := filesRedis.CollectFilesInFolder("/data")
		err := filesRedis.SaveFiles(files)
		assert.Nil(t, err)
	})

	t.Run("Should get file", func(t *testing.T) {
		file, err := filesRedis.GetFile("example1.log")
		assert.Nil(t, err)
		assert.NotEmpty(t, file)
		assert.Equal(t, "example1.log", file.Name)
		fullpath, _ := filepath.Abs("../data/example1.log")
		assert.Equal(t, fullpath, file.Path)
	})

	t.Run("Should cache content", func(t *testing.T) {
		err := filesRedis.CacheLogs("example1.log", []string{"line1", "line2"})
		assert.Nil(t, err)
	})
}
