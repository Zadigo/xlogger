package logic

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MetaData struct contains various boolean fields that indicate
// specific characteristics of the path of the request which can
// be useful for further vulnerability analysis.
type MetaData struct {
	// IsPHP indicates if the request is for a PHP file or contains "php" in the path
	IsPHP bool `json:"is_php"`
	// IsAssets indicates if the request is for a static asset (e.g., .css, .js, .png)
	IsAssets bool `json:"is_assets"`
	// IsJS indicates if the request is for a JavaScript file
	IsJS bool `json:"is_js"`
	// IsHTTP2 indicates if the request was made using the HTTP/2 protocol
	IsHTTP2 bool `json:"is_http2"`
	// IsRobotsTxt indicates if the request is for the robots.txt file
	IsRobotsTxt bool `json:"is_robots_txt"`
	// IsXml indicates if the request is for an XML file
	IsXml bool `json:"is_xml"`
	// IsAttemptedLogin indicates if the request is an attempted login
	// based on the presence of "login" in the path or user agent
	IsAttemptedLogin bool `json:"is_attempted_login"`
	// IsWordpress indicates if the request is related to a WordPress site
	IsWordpress bool `json:"is_wordpress"`
	// IsEnv indicates if the request is for an environment file
	IsEnv bool `json:"is_env"`
	// IsExecutable indicates if the request is for an executable file
	IsExecutable bool `json:"is_executable"`
	// IsPowerShell indicates if the request is for a PowerShell script
	IsPowerShell bool `json:"is_powershell"`
	// IsNuxt indicates if the request is related to a Nuxt.js application
	IsNuxt bool `json:"is_nuxt"`
	// IsGponRouter indicates if the request is related to a GPON router
	IsGponRouter bool `json:"is_gpon_router"`
}

type LogLine struct {
	RawLine string `json:"raw_line"`
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
	// MetaData contains various boolean fields that indicate specific
	// characteristics of the path of the request which can be useful
	// for further vulnerability analysis.
	MetaData MetaData `json:"meta_data"`
}

// Checks the value of the status code and returns
// if it was successful or not
func (l LogLine) analyzeStatusCode(status int) bool {
	return status >= 200 && status <= 226
}

// Parses a line of the log file and returns a LogLine struct
func (l LogLine) ParseLine() (LogLine, error) {
	// For ^(\S+) - (\S+) \[(.*)\] "(POST|GET|OPTIONS|PATCH|PUT) (.*) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "(\S+)" "(\S+)" (\d+) "(\S+)" "(\S+)" (\d+)ms$
	logLineRegex := regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH) ([^"]+) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`)

	var matched []string = logLineRegex.FindStringSubmatch(l.RawLine)
	if matched == nil {
		return LogLine{}, fmt.Errorf("🔴 Line is not valid %s", l.RawLine)
	}

	status, _ := strconv.Atoi(matched[7])
	l.IsSuccess = l.analyzeStatusCode(status)

	l.RemoteAddress = matched[1]
	l.RemoteUser = matched[2]
	l.DateTime = matched[3]
	l.Method = matched[4]
	l.Path = matched[5]
	l.Protocole = matched[6]
	l.StatusCode = status
	l.Referrer = matched[8]
	l.UserAgent = matched[9]
	l.BodyBytesSent, _ = strconv.Atoi(matched[7])

	// Date parsing
	dateLayout := "02/Jan/2006:15:04:05 -0700"
	parsedDate, err := time.Parse(dateLayout, matched[3])

	if err == nil {
		l.RemoteDate = parsedDate.Format("2006-01-02")
		l.RemoteTime = parsedDate.Format("15:04:05")
	}

	parsedPathExtension := filepath.Ext(l.Path)

	switch parsedPathExtension {
	case ".php":
		l.MetaData.IsPHP = true
	case ".js":
		l.MetaData.IsJS = true
	case ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico":
		l.MetaData.IsAssets = true
	case ".xml":
		l.MetaData.IsXml = true
	case ".env":
		l.MetaData.IsEnv = true
	case ".exe", ".sh", ".bat", ".cmd", ".ini", ".conf", ".config", ".bak", ".backup":
		l.MetaData.IsExecutable = true
	}

	if strings.Contains(l.Path, "php") || strings.Contains(l.Path, "laravel") {
		l.MetaData.IsPHP = true
	}

	if l.Path == "/robots.txt" {
		l.MetaData.IsRobotsTxt = true
	}

	if l.Protocole == "HTTP/2.0" {
		l.MetaData.IsHTTP2 = true
	}

	if strings.Contains(strings.ToLower(l.UserAgent), "login") || strings.Contains(strings.ToLower(l.Path), "login") {
		l.MetaData.IsAttemptedLogin = true
	}

	if strings.Contains(strings.ToLower(l.UserAgent), "wordpress") || strings.Contains(strings.ToLower(l.Path), "wp-") {
		l.MetaData.IsWordpress = true
	}

	if strings.Contains(l.Path, "cgi-bin") {
		l.MetaData.IsExecutable = true
	}

	if strings.Contains(strings.ToLower(l.Path), "powershell") {
		l.MetaData.IsPowerShell = true
	}

	if strings.Contains(l.Path, "_nuxt") {
		l.MetaData.IsNuxt = true
	}

	if strings.Contains(strings.ToLower(l.Path), "gpon") {
		l.MetaData.IsGponRouter = true
	}

	return l, nil
}
