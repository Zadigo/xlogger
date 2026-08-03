package tickerapp

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
	IsPHP bool `json:"isPhp"`
	// IsAssets indicates if the request is for a static asset (e.g., .css, .js, .png)
	IsAssets bool `json:"isAssets"`
	// IsJS indicates if the request is for a JavaScript file
	IsJS bool `json:"isJs"`
	// IsHTTP2 indicates if the request was made using the HTTP/2 protocol
	IsHTTP2 bool `json:"isHttp2"`
	// IsRobotsTxt indicates if the request is for the robots.txt file
	IsRobotsTxt bool `json:"isRobotsTxt"`
	// IsXml indicates if the request is for an XML file
	IsXml bool `json:"isXml"`
	// IsAttemptedLogin indicates if the request is an attempted login
	// based on the presence of "login" in the path or user agent
	IsAttemptedLogin bool `json:"isAttemptedLogin"`
	// IsWordpress indicates if the request is related to a WordPress site
	IsWordpress bool `json:"isWordpress"`
	// IsEnv indicates if the request is for an environment file
	IsEnv bool `json:"isEnv"`
	// IsExecutable indicates if the request is for an executable file
	IsExecutable bool `json:"isExecutable"`
	// IsPowerShell indicates if the request is for a PowerShell script
	IsPowerShell bool `json:"isPowershell"`
	// IsNuxt indicates if the request is related to a Nuxt.js application
	IsNuxt bool `json:"isNuxt"`
	// IsGponRouter indicates if the request is related to a GPON router
	IsGponRouter bool `json:"isGponRouter"`
	// IsWindows indicates the request tried to access a Windows path (e.g., contains backslashes)
	IsWindowsPath bool `json:"isWindows"`
	// IsGitHub indicates if the request is related to GitHub (e.g., contains ".git" in the path or user agent)
	IsGitHub bool `json:"isGithub"`
}

type LogLine struct {
	RawLine string `json:"rawline"`
	// The IP address of the client
	RemoteAddress string `json:"remoteAddress"`
	// The authenticated user
	RemoteUser string `json:"remoteUser"`
	// Date and time at which the request was made + TZ
	DateTime string `json:"datetime"`
	// Method used for the request
	Method string `json:"method"`
	// The path of the request
	Path string `json:"path"`
	// The HTTP protocole used e.g HTTP/2.0
	Protocole string `json:"protocole"`
	// The request status code
	StatusCode int `json:"statusCode"`
	// Number of bytes sent to the client (body only, not headers)
	BodyBytesSent int `json:"bodyBytesSent"`
	// The page from which the user came
	Referrer string `json:"referrer"`
	// The client's user agent
	UserAgent string `json:"user_agent"`
	// The date part of the date time
	RemoteDate string `json:"remoteDate"`
	// the time part of the date time
	RemoteTime string `json:"remoteTime"`
	// Whether the request was successful
	IsSuccess bool `json:"isSuccess"`
	// MetaData contains various boolean fields that indicate specific
	// characteristics of the path of the request which can be useful
	// for further vulnerability analysis.
	MetaData MetaData `json:"metaData"`
}

// Checks the value of the status code and returns
// if it was successful or not
func (l *LogLine) analyzeStatusCode(status int) bool {
	return status >= 200 && status <= 226
}

// Parses a line of the log file and returns a LogLine struct
func (l *LogLine) ParseLine() (line *LogLine, err error) {
	// For ^(\S+) - (\S+) \[(.*)\] "(POST|GET|OPTIONS|PATCH|PUT) (.*) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "(\S+)" "(\S+)" (\d+) "(\S+)" "(\S+)" (\d+)ms$
	logLineRegex := regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH) ([^"]+) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`)

	var matched []string = logLineRegex.FindStringSubmatch(l.RawLine)
	if matched == nil {
		return &LogLine{}, fmt.Errorf("🔴 Line is not valid %s", l.RawLine)
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
	l.Referrer = matched[9]
	l.UserAgent = matched[10]
	l.BodyBytesSent, _ = strconv.Atoi(matched[8])

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

	if strings.Contains(l.Path, "/.git") {
		l.MetaData.IsGitHub = true
	}

	regex := regexp.MustCompile(`\w+\:\/[Ww]indows`)
	matched = regex.FindStringSubmatch(l.Path)
	if matched != nil {
		l.MetaData.IsWindowsPath = true
	}

	return l, nil
}
