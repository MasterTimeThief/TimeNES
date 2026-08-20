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
	PendingMute       bool
	PendingScreenshot bool
	PendingFullscreen bool
)

var (
	UIMessage      string
	UIMessageTimer int
)
