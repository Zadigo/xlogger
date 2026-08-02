package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestLogServerConfig(t *testing.T) {
	config := utils.ServerConfig{}
	t.Run("Should load YAML config", func(t *testing.T) {
		err := config.Load(t.Context())
		assert.Nil(t, err)
		assert.Equal(t, "* * * * *", config.LogServer.Interval)
		assert.Equal(t, "data", config.LogServer.Logs.Folder)
	})
}
