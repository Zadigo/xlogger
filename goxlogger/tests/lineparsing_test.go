package tests

import (
	"fmt"
	"testing"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	type TestCase struct {
		name               string
		line               string
		shouldParse        bool
		expectedMethod     string
		expectedStatusCode int
		expectedPHP        bool
		expectedEnv        bool
	}

	testCases := []TestCase{
		{
			name:               "Simple GET request",
			line:               `1.1.1.1 - - [08/May/2025:13:53:52 +0000] "GET / HTTP/1.1" 200 701 "https://example.com" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`,
			shouldParse:        true,
			expectedMethod:     "GET",
			expectedStatusCode: 200,
			expectedPHP:        false,
			expectedEnv:        false,
		},
		{
			name:               "Simple POST request",
			line:               `1.1.1.1 - - [08/May/2025:13:53:52 +0000] "POST / HTTP/1.1" 200 701 "https://example.com" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`,
			shouldParse:        true,
			expectedMethod:     "POST",
			expectedStatusCode: 200,
			expectedPHP:        false,
			expectedEnv:        false,
		},
		{
			name:               "PHP file request",
			line:               `1.1.1.1 - - [13/Nov/2025:03:46:34 +0000] "GET /example.php HTTP/1.1" 200 312 "-" "-"`,
			shouldParse:        true,
			expectedMethod:     "GET",
			expectedStatusCode: 200,
			expectedPHP:        true,
			expectedEnv:        false,
		},
		{
			name:               "Env file request",
			line:               `1.1.1.1 - - [12/Nov/2025:12:31:10 +0000] "GET /.env HTTP/1.1" 200 312 "-" "Python/3.10"`,
			shouldParse:        true,
			expectedMethod:     "GET",
			expectedStatusCode: 200,
			expectedPHP:        false,
			expectedEnv:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawLine := logic.LogLine{RawLine: tc.line}
			line, err := rawLine.ParseLine()

			assert.NotEmpty(t, line)
			assert.NoError(t, err)

			if tc.shouldParse {
				assert.NotEmpty(t, line.Method, fmt.Sprintf("Method should not be empty: %s", line.Method))
				assert.Equal(t, tc.expectedMethod, line.Method, fmt.Sprintf("Expected method %s, got %s", tc.expectedMethod, line.Method))
				assert.Equal(t, tc.expectedStatusCode, line.StatusCode, fmt.Sprintf("Expected status code %d, got %d", tc.expectedStatusCode, line.StatusCode))
				assert.Equal(t, "1.1.1.1", line.RemoteAddress, fmt.Sprintf("Expected remote address %s, got %s", "1.1.1.1", line.RemoteAddress))

				if tc.expectedPHP {
					assert.Contains(t, line.Path, ".php", fmt.Sprintf("Expected path to contain .php, got %s", line.Path))
					assert.True(t, line.MetaData.IsPHP, "Expected IsPHP to be true")
				}
			}
		})
	}
}
