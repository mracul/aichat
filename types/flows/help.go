package flows
// help.go
// Type definitions for help-related flows (show controls, show about)

package flows

import "aichat/types"

// ShowControlsFlow for displaying controls cheat sheet
type ShowControlsFlow struct {
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ShowAboutFlow for displaying about information
type ShowAboutFlow struct {
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

