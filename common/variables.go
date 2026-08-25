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
)

var (
	UIMessage      string
	UIMessageTimer int
)

func Reset() {
	CPU_TotalCycles = 0
	UIMessage = ""
	UIMessageTimer = 0
}
