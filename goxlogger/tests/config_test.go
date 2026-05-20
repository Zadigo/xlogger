package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestLogServerConfig(t *testing.T) {
	config := models.YamlConfig{}
	t.Run("Should load YAML config", func(t *testing.T) {
		err := config.Load("..")
		assert.Nil(t, err)
		assert.Equal(t, "* * * * *", config.LogServer.Interval)
		assert.Equal(t, "data", config.LogServer.Logs.Folder)
	})
}
