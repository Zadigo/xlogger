package test

import (
	"testing"

	"github.com/Zadigo/xlogger_backend/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestServerConfig(t *testing.T) {
	config := backend.GetServerConfig()
	t.Run("Load the configuration file", func(t *testing.T) {
		assert.NotNil(t, config.Config.LogsFolder.Name)
	})
}
