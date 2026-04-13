package backend

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type LogLine struct {
	// The IP address of the client
	RemoteAddress string `json:"remote_address"`
	// The authenticated user
	RemoteUser string `json:"remote_user"`
	// Date and time at which the request was made + TZ
	DateTime string `json:"date_time"`
	// Method used for the request
	Method string `json:"method"`
	// The path of the request
	Path string `json:"path"`
	// The HTTP protocole used e.g HTTP/2.0
	Protocole string `json:"protocole"`
	// The request status code
	StatusCode int `json:"status_code"`
	// Number of bytes sent to the client (body only, not headers)
	BodyBytesSent int `json:"body_bytes_sent"`
	// The page from which the user came
	Referrer string `json:"referrer"`
	// The client's user agent
	UserAgent string `json:"user_agent"`
	// The date part of the date time
	RemoteDate string `json:"remote_date"`
	// the time part of the date time
	RemoteTime string `json:"remote_time"`
	// Whether the request was successful
	IsSuccess bool `json:"is_success"`
}

// Checks the value of the status code and returns
// if it was successful or not
func AnalyzeStatusCode(status int) bool {
	return status >= 200 && status <= 226
}

func AnalyzePath(logLine LogLine) string {
	switch logLine.Path {
	case "/":
		return "home"
	}
	return "unknown"
}

// Parses a line of the log file and returns a LogLine struct
func ParseLine(textValue string) (LogLine, error) {
	logLineRegex := regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH) ([^"]+) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`)

	var matched []string = logLineRegex.FindStringSubmatch(textValue)
	if matched == nil {
		return LogLine{}, fmt.Errorf("🔴 Line is not valid %s", textValue)
	}

	status, _ := strconv.Atoi(matched[7])
	bytes, _ := strconv.Atoi(matched[8])

	// Parse the actual date and time from
	// the datetime string that we got

	var remoteDate string
	var remoteTime string

	dateLayout := "02/Jan/2006:15:04:05 -0700"
	parsedDate, err := time.Parse(dateLayout, matched[3])

	if err == nil {
		remoteDate = parsedDate.Format("2006-01-02")
		remoteTime = parsedDate.Format("15:04:05")
	}

	logLine := LogLine{
		RemoteAddress: matched[1],
		RemoteUser:    matched[2],
		DateTime:      matched[3],
		Method:        matched[4],
		Path:          matched[5],
		Protocole:     matched[6],
		StatusCode:    status,
		BodyBytesSent: bytes,
		Referrer:      matched[9],
		UserAgent:     matched[10],

		RemoteDate: remoteDate,
		RemoteTime: remoteTime,

		IsSuccess: AnalyzeStatusCode(status),
	}

	return logLine, nil
}
