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

func TestAndroidAnalyzer(t *testing.T) {
	testlines := []string{
		"Mozilla/5.0 (Linux; Android 4.4.3; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/34.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY48I) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.84 Safari/537.36",
	}

	for _, line := range testlines {
		t.Run(fmt.Sprintf("Testing line: %s", line), func(t *testing.T) {
			agent := &logic.UserAgent{RawValue: line}

			analyzer := &logic.AndroidAnalyzer{}
			analyzer.Execute(agent)

			assert.Equal(t, "Linux;", agent.OperatingSystemPlatform, fmt.Sprintf("Expected platform 'Linux;', got '%s'", agent.OperatingSystemPlatform))
			assert.Equal(t, "Android", agent.OperatingSystemKernel, fmt.Sprintf("Expected OS 'Android', got '%s'", agent.OperatingSystemKernel))
			assert.NotEmpty(t, agent.OperatingSystemBuildVersion, "Expected OS version to be non-empty")
		})
	}
}

func TestMacAnalyzer(t *testing.T) {
	testlines := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/600.4.10 (KHTML, like Gecko) Version/7.1.4 Safari/537.85.13",
	}

	for _, line := range testlines {
		t.Run(fmt.Sprintf("Testing line: %s", line), func(t *testing.T) {
			agent := &logic.UserAgent{RawValue: line}

			analyzer := &logic.MacAnalyzer{}
			analyzer.Execute(agent)

			assert.Equal(t, "Macintosh", agent.OperatingSystemPlatform, fmt.Sprintf("Expected platform 'Macintosh', got '%s'", agent.OperatingSystemPlatform))
			assert.Equal(t, "Intel", agent.OperatingSystemKernel, fmt.Sprintf("Expected OS 'Intel', got '%s'", agent.OperatingSystemKernel))
			assert.NotEmpty(t, agent.OperatingSystemBuildVersion, "Expected OS version to be non-empty")
		})
	}
}

func TestWindowsAnalyzer(t *testing.T) {
	type TestCase struct {
		name string
		line string
	}

	testCases := []TestCase{
		{name: "test_with_wow", line: "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:38.0) Gecko/20100101 Firefox/38.0"},
		{name: "test_with_trident", line: "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko"},
		{name: "test_with_touch", line: "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; ASU2JS; rv:11.0) like Gecko"},
		{name: "test_with_touch", line: "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; MALNJS; rv:11.0) like Gecko"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			agent := &logic.UserAgent{RawValue: tc.line}

			analyzer := &logic.WindowsAnalyzer{}
			analyzer.Execute(agent)

			assert.Equal(t, "Windows NT", agent.OperatingSystemPlatform, fmt.Sprintf("Expected platform 'Windows', got '%s'", agent.OperatingSystemPlatform))
			assert.NotEmpty(t, agent.OperatingSystemBuildVersion, "Expected OS version to be non-empty")
		})
	}
}
