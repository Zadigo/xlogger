package tests

import (
	"encoding/json"
	"testing"

	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	recorder := CreateGetFilesRecorder(t)

	t.Run("should return files", func(t *testing.T) {
		var files []tickerapp.File
		err := json.Unmarshal(recorder.Body.Bytes(), &files)

		assert.NoError(t, err)
		assert.Equal(t, recorder.Code, 200)
	})
}

func TestGetLogs(t *testing.T) {
	recorder := CreateGetLogsRecorder(t)

	t.Run("should return logs", func(t *testing.T) {
		var logs []tickerapp.LogLine
		err := json.Unmarshal(recorder.Body.Bytes(), &logs)

		assert.NoError(t, err)
		assert.Equal(t, recorder.Code, 200)
	})
}
