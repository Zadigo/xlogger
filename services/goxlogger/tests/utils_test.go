package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	t.Run("should get root dir when not in context", func(t *testing.T) {
		r := utils.GetRootDir(t.Context())
		assert.NotEmpty(t, r)
		assert.True(t, strings.Contains(r, "/test"))
	})

	t.Run("should get root dir when not in context", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), "rootDir", "..")
		r := utils.GetRootDir(ctx)
		assert.NotEmpty(t, r)
		assert.False(t, strings.Contains(r, "/test"))
	})
}
