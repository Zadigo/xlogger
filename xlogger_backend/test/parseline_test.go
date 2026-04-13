package test

import (
	"testing"

	"github.com/Zadigo/xlogger_backend/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	testLines := []string{
		"172.21.0.2 - - [08/May/2025:13:53:43 +0000] \"GET /assets/index-Dl4nsTUV.css HTTP/1.1\" 200 9274 \"https://example.fr/\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36\"",
		"172.21.0.2 - - [08/May/2025:13:53:43 +0000] \"GET /assets/DefaultSite-BiGz7DoW.js HTTP/1.1\" 200 208 \"-\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36\"",
	}

	type testCase struct {
		testLine string
		expected bool
	}

	testCases := []testCase{
		{
			testLine: testLines[0],
			expected: true,
		},
		{
			testLine: testLines[1],
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run("Extract simple line", func(t *testing.T) {
			result, err := backend.LogLine{RawLine: tc.testLine}.ParseLine()
			assert.NoError(t, err)
			assert.NotNil(t, result.RemoteDate)
			assert.NotNil(t, result.RemoteAddress)
		})
	}
}
