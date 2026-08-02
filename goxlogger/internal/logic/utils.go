package logic

import (
	"regexp"
	"strings"
)

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

func (a *UserAgent) MatchRegex(regex string) []string {
	re := regexp.MustCompile(regex)
	return re.FindStringSubmatch(a.RawValue)
}

// Analyzes the user agent string and returns a UserAgent struct
// func (u *UserAgent) AnalyzeUserAgent(value string) (agent UserAgent) {
// 	browserRegex := u.MatchRegex(`^((Chrome|Mozilla)\/(\d+\.\d+))`, value)
// 	if browserRegex != nil {
// 		agent.Name = browserRegex[2]
// 		agent.Version = browserRegex[3]
// 	}

// 	matchedMac := u.MatchRegex(`(?:.*)(Macintosh)\; (Intel)? Mac (OS \w+) (.*\d+)\)\s`, value)
// 	if matchedMac != nil {
// 		agent.Platform = matchedMac[1]
// 		agent.OS = matchedMac[3]
// 		agent.OsVersion = matchedMac[4]
// 	}

// 	matchedAndroid := u.MatchRegex(`(?:.*)(Linux(?:\;)) (Android)? (\d+\.\d+\.?\d?); (.*)\)\s?(?=AppleWebKit)`, value)
// 	if matchedAndroid != nil {
// 		agent.Platform = matchedAndroid[1]
// 		agent.OS = matchedAndroid[2]
// 		agent.OsVersion = matchedAndroid[3]
// 	}

// 	return agent
// }

type AnalyzerInterface interface {
	execute(agent *UserAgent)
	setNext(AnalyzerInterface)
}

type BaseAnalyzer struct {
	next AnalyzerInterface
}

func (b *BaseAnalyzer) setNext(next AnalyzerInterface) {
	b.next = next
}

func (b *BaseAnalyzer) execute(agent *UserAgent) {
	if b.next != nil {
		b.next.execute(agent)
	}
}

type AndroidAnalyzer struct {
	next AnalyzerInterface
}

func (a *AndroidAnalyzer) execute(agent *UserAgent) {
	if strings.Contains(agent.RawValue, "Android") {
		matchedAndroid := agent.MatchRegex(`(?:.*)(Linux(?:\;)) (Android)? (\d+\.\d+\.?\d?); (.*)\)\s?(?=AppleWebKit)`)
		if matchedAndroid != nil {
			agent.Platform = matchedAndroid[1]
			agent.OS = matchedAndroid[2]
			agent.OsVersion = matchedAndroid[3]
		}
	} else {
		if a.next != nil {
			a.next.execute(agent)
		}
	}
}

func (a *AndroidAnalyzer) setNext(next AnalyzerInterface) {
	a.next = next
}

type MacAnalyzer struct {
	next AnalyzerInterface
}

func (m *MacAnalyzer) execute(agent *UserAgent) {
	if strings.Contains(agent.RawValue, "Macintosh") || strings.Contains(agent.RawValue, "Mac OS") {
		matchedMac := agent.MatchRegex(`(?:.*)(Macintosh)\; (Intel)? Mac (OS \w+) (.*\d+)\)\s`)
		if matchedMac != nil {
			agent.Platform = matchedMac[1]
			agent.OS = matchedMac[3]
			agent.OsVersion = matchedMac[4]
		}
	} else {
		if m.next != nil {
			m.next.execute(agent)
		}
	}
}

func (m *MacAnalyzer) setNext(next AnalyzerInterface) {
	m.next = next
}

type BrowserAnalyzer struct {
	next AnalyzerInterface
}

func (b *BrowserAnalyzer) execute(agent *UserAgent) {
	matched := agent.MatchRegex(`^((Chrome|Mozilla)\/(\d+\.\d+))`)
	if matched != nil {
		agent.Name = matched[2]
		agent.Version = matched[3]
	} else {
		if b.next != nil {
			b.next.execute(agent)
		}
	}
}

func (b *BrowserAnalyzer) setNext(next AnalyzerInterface) {
	b.next = next
}

func NewUserAgentAnalysis(value string) *UserAgent {
	agent := &UserAgent{RawValue: value}

	baseAnalyzer := &BaseAnalyzer{}

	macAnalyzer := &MacAnalyzer{}
	macAnalyzer.setNext(baseAnalyzer)

	androidAnalyzer := &AndroidAnalyzer{}
	androidAnalyzer.setNext(macAnalyzer)

	browserAnalyzer := &BrowserAnalyzer{}
	browserAnalyzer.setNext(androidAnalyzer)

	browserAnalyzer.execute(agent)
	return agent
}
