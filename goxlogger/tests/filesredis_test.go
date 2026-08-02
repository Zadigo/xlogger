package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Zadigo/goxlogger/internal/backend"
	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/stretchr/testify/assert"
)

func instanceFixture() *logic.FileRedis {
	ctx := context.WithValue(context.Background(), "rootDir", "../")

	redisClient := backend.NewRedisBackend()
	filesRedis := logic.NewFileRedis(ctx, redisClient)

	return filesRedis
}

func TestGetLocalLogs(t *testing.T) {
	filesRedis := instanceFixture()
	files, err := filesRedis.GetLocalLogs("/data")

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

func TestImplementation(t *testing.T) {
	ctx := context.WithValue(t.Context(), "rootDir", "../")

	redisClient := backend.NewRedisBackend()
	filesRedis := logic.NewFileRedis(ctx, redisClient)

	t.Run("Should load files", func(t *testing.T) {
		files, err := filesRedis.GetLocalLogs("/data")
		assert.Nil(t, err)
		assert.NotEmpty(t, files)
		assert.NoError(t, err)
		assert.True(t, len(files) > 0)
	})

	t.Run("Should save files", func(t *testing.T) {
		files, _ := filesRedis.GetLocalLogs("/data")
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
		err := filesRedis.CacheContent("example1.log", []string{"line1", "line2"})
		assert.Nil(t, err)
	})
}
