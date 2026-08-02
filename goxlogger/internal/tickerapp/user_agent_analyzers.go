package tickerapp

import (
	"log"
	"regexp"
	"strings"
)

type Architecture struct {
	Version string `json:"version"`
	Is64Bit bool   `json:"is_64_bit"` // Indicates whether the user agent is running on a 64-bit architecture

}

type Windows struct {
	Version              string `json:"version"`
	LayoutEngine         string `json:"layout_engine"`
	LayoutEngineVersion  string `json:"layout_engine_version"`
	BrowserVersion       string `json:"browser_version"`
	HasCompatibilityView bool   `json:"has_compatibility_view"` // Indicates whether the browser is in compatibility view mode
	MSIEVersion          string `json:"msie_version"`           // Microsoft Internet Explorer (MSIE) version information, if applicable
	IsMSIE               bool   `json:"is_msie"`                // Indicates whether the user agent is Microsoft Internet Explorer (MSIE)
	Is32Bit              bool   `json:"is_32_bit"`              // Indicates whether the user agent is running on a 32-bit architecture
	Is64Bit              bool   `json:"is_64_bit"`              // Indicates whether the user agent is running on a 64-bit architecture
	IsTouch              bool   `json:"is_touch"`               // Indicates whether the device has touch capabilities
}

type Macintosh struct {
	Architecture
	ProcessorArchitecture string `json:"processor_architecture"`
	IsIntel               bool   `json:"is_intel"`   // Indicates whether the user agent is running on an Intel-based Macintosh
	IsSiliconChip         bool   `json:"is_silicon"` // Indicates whether the user agent is running on an Apple Silicon-based Macintosh
}

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
	// The WindowsInfo field contains additional information about the Windows operating system, if applicable
	WindowsInfo *Windows `json:"windows_info,omitempty"`
	// The MacintoshInfo field contains additional information about the Macintosh operating system, if applicable
	MacintoshInfo *Macintosh `json:"macintosh_info,omitempty"`
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
		agent.OperatingSystemPlatform = "Macintosh"
		agent.MacintoshInfo = &Macintosh{}

		matched := agent.MatchRegex(`Macintosh; (Intel)`)
		if matched != nil {
			agent.MacintoshInfo.IsIntel = true
		}

		matched = agent.MatchRegex(`Macintosh;.*Mac OS (\w+) (\d+\_\d+\_\d+)`)
		if matched != nil {
			agent.OperatingSystemBuildVersion = strings.ReplaceAll(matched[2], "_", ".")
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

// WindowsAnalyzer is responsible for analyzing user agent strings that
// contain "Windows" and extracting relevant information.
type WindowsAnalyzer struct {
	next AnalyzerInterface
}

func (w *WindowsAnalyzer) Execute(agent *UserAgent) {
	if strings.Contains(agent.RawValue, "Windows") {
		agent.OperatingSystemPlatform = "Windows NT"

		// e.g. (Windows NT 6.3; WOW64)
		matched := agent.MatchRegex(`Windows NT\s(\d+\.\d?); (WOW\d+)`)
		if matched != nil {
			agent.OperatingSystemBuildVersion = matched[1]
			agent.WindowsInfo = &Windows{
				Version: matched[1],
				Is32Bit: matched[2] == "WOW32",
				Is64Bit: matched[2] == "WOW64",
			}
		}

		// e.g. (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0)
		matched = agent.MatchRegex(`Windows NT.*(Trident\W(\d+\.\d+))`)
		if matched != nil {
			agent.WindowsInfo.LayoutEngine = matched[1]
			agent.WindowsInfo.LayoutEngineVersion = matched[2]
		}

		matched = agent.MatchRegex(`Windows NT.*(Touch)`)
		if matched != nil {
			agent.WindowsInfo.IsTouch = true
		}

		matched = agent.MatchRegex(`Windows NT.*(MSIE\W(\d+\.\d+))`)
		if matched != nil {
			agent.WindowsInfo.MSIEVersion = matched[2]
			agent.WindowsInfo.IsMSIE = true
			agent.WindowsInfo.BrowserVersion = matched[2]
		}

		// e.g. (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0)
		matched = agent.MatchRegex(`Windows NT.*(rv\:(\d+\.\d+))`)
		if matched != nil {
			agent.WindowsInfo.BrowserVersion = matched[2]
		}
	} else {
		if w.next != nil {
			w.next.Execute(agent)
		}
	}
}

func (w *WindowsAnalyzer) setNext(next AnalyzerInterface) {
	w.next = next
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
	agent := &UserAgent{
		RawValue:    value,
		WindowsInfo: &Windows{},
	}

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
