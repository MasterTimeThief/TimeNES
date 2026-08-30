package common

var CPU_TotalCycles int

var (
	ScreenScale    = 2
	NewScreenScale = -1
)

var (
	PendingReset      bool
	PendingPause      bool
	PendingFPS        bool
	PendingDebug      bool
	PendingPT         bool
	PendingMute       bool
	PendingScreenshot bool
	PendingFullscreen bool

	PendingPatternUpdate bool

	UIMessage      string
	UIMessageTimer int

	EmulatorVolume float64 = 25.0 // out of 100

	Toggle90High  bool
	Toggle440High bool
	Toggle14KLow  bool
)

func Reset() {
	CPU_TotalCycles = 0
	UIMessage = ""
	UIMessageTimer = 0
}
