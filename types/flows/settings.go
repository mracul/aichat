// settings.go
// Type definitions for settings-related flows

package flows

import "aichat/types"

// SettingsFlow for multi-step settings/preferences modal
type SettingsFlow struct {
	ThemeSelector types.ViewState
	KeybindEditor types.ViewState
	NoticeModal   types.ViewState
	Strategy      FlowStrategy
}
