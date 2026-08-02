package logic

import (
	"log"
	"regexp"
	"strings"
)

type UserAgent struct {
	// The raw user agent string as received from the client
	RawValue string `json:"raw_value"`
	// The version of the browser or client making the request. It is a
	// legacy identifier used by almost all modern browsers to ensure
	// web compatibility with web servers
	ProductToken string `json:"product_token"`
	// The version of the browser or client making the request
	ProductVersion string `json:"version"`
	// The operating system kernel of the platform (e.g., Linux, Darwin) which
	// indicates the base platform architecture
	OperatingSystemKernel string `json:"operating_system_kernel"`
	// The operating system platform (e.g., Windows, macOS, Linux) which indicates the general platform type
	OperatingSystemPlatform string `json:"platform"`
	// The hardware model identifier of the platform (e.g., iPhone, iPad, MacBookPro)
	// which indicates the specific device model
	HardwareModelIdentifier string `json:"hardware_model_identifier"`
	// The operating system version of the platform (e.g., 10.15.7, 11.2.3)
	OperatingSystemBuildVersion string `json:"operating_system_build_version"`
}

func (a *UserAgent) MatchRegex(regex string) []string {
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Print(err)
		return nil
	}
	return re.FindStringSubmatch(a.RawValue)
}

type AnalyzerInterface interface {
	// Execute analyzes the given UserAgent and extracts relevant information.
	Execute(agent *UserAgent)
	// setNext sets the next analyzer in the chain.
	setNext(AnalyzerInterface)
}

// BaseAnalyzer is a struct that implements the
// AnalyzerInterface and serves as a base for other analyzers.
type BaseAnalyzer struct {
	next AnalyzerInterface
}

func (b *BaseAnalyzer) setNext(next AnalyzerInterface) {
	b.next = next
}

func (b *BaseAnalyzer) Execute(agent *UserAgent) {
	if b.next != nil {
		b.next.Execute(agent)
	}
}

// AndroidAnalyzer is responsible for analyzing user agent strings that
// contain "Android" and extracting relevant information.
type AndroidAnalyzer struct {
	next AnalyzerInterface
}

func (a *AndroidAnalyzer) Execute(agent *UserAgent) {
	if strings.Contains(agent.RawValue, "Android") {
		matchedAndroid := agent.MatchRegex(`(?:.*)(Linux(?:\;)) (Android)? (\d+\.\d+\.?\d?); (.*)\)\s?AppleWebKit`)
		if matchedAndroid != nil {
			agent.OperatingSystemPlatform = matchedAndroid[1]
			agent.OperatingSystemKernel = matchedAndroid[2]
			agent.OperatingSystemBuildVersion = matchedAndroid[3]
		}
	} else {
		if a.next != nil {
			a.next.Execute(agent)
		}
	}
}

func (a *AndroidAnalyzer) setNext(next AnalyzerInterface) {
	a.next = next
}

// MacAnalyzer is responsible for analyzing user agent strings that
// contain "Macintosh" or "Mac OS" and extracting relevant information.
type MacAnalyzer struct {
	next AnalyzerInterface
}

func (m *MacAnalyzer) Execute(agent *UserAgent) {
	if strings.Contains(agent.RawValue, "Macintosh") || strings.Contains(agent.RawValue, "Mac OS") {
		matchedMac := agent.MatchRegex(`(?:.*)(Macintosh)\; (Intel)? Mac (OS \w+) (.*\d+)\)\s`)
		if matchedMac != nil {
			agent.OperatingSystemPlatform = matchedMac[1]
			agent.OperatingSystemKernel = matchedMac[2]
			agent.OperatingSystemBuildVersion = matchedMac[4]
		}
	} else {
		if m.next != nil {
			m.next.Execute(agent)
		}
	}
}

func (m *MacAnalyzer) setNext(next AnalyzerInterface) {
	m.next = next
}

// BrowserAnalyzer is responsible for analyzing user agent strings that
// contain browser information (e.g., Chrome, Mozilla) and extracting relevant information.
type BrowserAnalyzer struct {
	next AnalyzerInterface
}

func (b *BrowserAnalyzer) Execute(agent *UserAgent) {
	matched := agent.MatchRegex(`^((Chrome|Mozilla)\/(\d+\.\d+))`)
	if matched != nil {
		agent.ProductToken = matched[2]
		agent.ProductVersion = matched[3]
	} else {
		if b.next != nil {
			b.next.Execute(agent)
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

	browserAnalyzer.Execute(agent)
	return agent
}
