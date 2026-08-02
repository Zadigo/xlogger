package logic

import "regexp"

type UserAgent struct {
	// The name of the browser or client making the request
	RawValue string `json:"raw_value"`
	// The name of the browser or client making the request
	Name string `json:"name"`
	// The version of the browser or client making the request
	Version string `json:"version"`
	// The platform on which the browser or client is running (e.g., Windows, macOS, Linux)
	Platform string `json:"platform"`
	// The operating system of the platform (e.g., Windows 10, macOS 11.2)
	OS string `json:"os"`
	// The version of the operating system
	OsVersion string `json:"os_version"`
}

func (u *UserAgent) MatchRegex(regex, value string) []string {
	re := regexp.MustCompile(regex)
	return re.FindStringSubmatch(value)
}

// Analyzes the user agent string and returns a UserAgent struct
func (u *UserAgent) AnalyzeUserAgent(value string) (agent UserAgent) {
	browserRegex := u.MatchRegex(`^((Chrome|Mozilla)\/(\d+\.\d+))`, value)
	if browserRegex != nil {
		agent.Name = browserRegex[2]
		agent.Version = browserRegex[3]
	}

	matchedMac := u.MatchRegex(`(?:.*)(Macintosh)\; (Intel)? Mac (OS \w+) (.*\d+)\)\s`, value)
	if matchedMac != nil {
		agent.Platform = matchedMac[1]
		agent.OS = matchedMac[3]
		agent.OsVersion = matchedMac[4]
	}

	matchedAndroid := u.MatchRegex(`(?:.*)(Linux(?:\;)) (Android)? (\d+\.\d+\.?\d?); (.*)\)\s?(?=AppleWebKit)`, value)
	if matchedAndroid != nil {
		agent.Platform = matchedAndroid[1]
		agent.OS = matchedAndroid[2]
		agent.OsVersion = matchedAndroid[3]
	}

	return agent
}
