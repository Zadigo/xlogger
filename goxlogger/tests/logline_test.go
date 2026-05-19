package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/stretchr/testify/assert"
)

func TestLogLine(t *testing.T) {
	type testCase struct {
		name  string
		line  string
		isPhp bool
	}

	testCases := []testCase{
		// {name: "favicon.ico request", line: `168.76.20.229 - - [20/Oct/2025:20:50:46 +0000] "GET /favicon.ico HTTP/1.1" 404 19 "-" "-" 6 "-" "-" 0ms`},
		{
			name:  "home page request",
			line:  `172.21.0.2 - - [08/May/2025:13:53:43 +0000] "GET / HTTP/1.1" 200 312 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`,
			isPhp: false,
		},
		{
			name:  "home page request",
			line:  `196.251.115.128 - - [20/Oct/2025:22:24:02 +0000] "POST /wp-confiq.php HTTP/1.1" 404 19 "-" "-" 5 "-" "-" 0ms`,
			isPhp: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line := logic.LogLine{RawLine: tc.line}
			parsedLine, err := line.ParseLine()
			assert.Nil(t, err)
			assert.Equal(t, "GET", parsedLine.Method)
			assert.Equal(t, 200, parsedLine.StatusCode)
			assert.False(t, parsedLine.MetaData.IsHTTP2)
			assert.True(t, parsedLine.IsSuccess)
			assert.Equal(t, tc.isPhp, parsedLine.MetaData.IsPHP)
		})
	}
}
