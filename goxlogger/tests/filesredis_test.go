package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestTestFilesRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	filesRedis := logic.NewFileRedis(t.Context(), "../data", redisClient)

	t.Run("Should load files", func(t *testing.T) {
		files, err := filesRedis.GetLocalLogs("/data")
		assert.Nil(t, err)
		assert.NotEmpty(t, files)
		fmt.Print(files)
	})

	t.Run("Should get files details", func(t *testing.T) {
		file := filesRedis.FileFromString("../data/test.log")
		assert.Equal(t, "test.log", file.Name)
		assert.Equal(t, "../data/test.log", file.Path)
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
