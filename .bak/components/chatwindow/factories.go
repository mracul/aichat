package chatwindow
package chatwindow

import (
	"aichat/components/modals"
	"aichat/types"
	"aichat/types/render"

	tea "github.com/charmbracelet/bubbletea"
)

// Factory for ModalViewState
func NewModalViewStateFactory(prev *ChatWindowViewState, content string, onYes, onNo func() tea.Msg, manager *modals.ModalManager, themeMap render.ThemeMap, strategy render.RenderStrategy) *types.ModalViewState {
	return &types.ModalViewState{
		ModalType: "chatwindow_modal",
		Content:   content,
	}
}

// Factory for ChatWindowViewState
func NewChatWindowViewStateFactory(chatID string, messages []types.Message, metadata types.ChatMetadata, inputBuffer, focus string, themeMap render.ThemeMap, strategy render.RenderStrategy) *ChatWindowViewState {
	return &ChatWindowViewState{
		ChatID:         chatID,
		Messages:       messages,
		Metadata:       metadata,
		InputBuffer:    inputBuffer,
		Focus:          focus,
		ThemeMap:       themeMap,
		RenderStrategy: strategy,
	}
}
