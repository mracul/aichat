// chats.go
// Type definitions for chat-related flows (new, custom, load, delete, list, export, rename)

package flows

import "aichat/types"

// FlowStrategy defines the interface for step/strategy logic in a flow
// Used to render and handle each modal step in a flow
// (Strategy pattern for flow steps)
type FlowStrategy interface {
	NextStep() types.ViewState
	PrevStep() types.ViewState
	CurrentStep() types.ViewState
	IsComplete() bool
}

// NewChatFlow defines the structure for the new chat creation flow
// Includes input modals for title, prompt, etc.
type NewChatFlow struct {
	TitleInput  types.ViewState // Input modal for chat title
	PromptInput types.ViewState // Input modal for prompt
	Strategy    FlowStrategy
}

// CustomChatFlow for advanced chat creation (name, model, prompt, system message)
type CustomChatFlow struct {
	NameInput     types.ViewState
	ModelSelector types.ViewState
	PromptInput   types.ViewState
	SystemInput   types.ViewState
	Strategy      FlowStrategy
}

// LoadChatFlow for selecting and loading a saved chat
type LoadChatFlow struct {
	ListModal   types.ViewState // List modal for saved chats
	NoticeModal types.ViewState // Notice modal for errors/info
	Strategy    FlowStrategy
}

// DeleteChatFlow for confirming and deleting a chat
type DeleteChatFlow struct {
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// ListChatsFlow for listing and selecting chats
type ListChatsFlow struct {
	ListModal types.ViewState
	Strategy  FlowStrategy
}

// ExportChatFlow for exporting a chat transcript
type ExportChatFlow struct {
	FormatSelector types.ViewState
	NoticeModal    types.ViewState
	Strategy       FlowStrategy
}

// RenameChatFlow for renaming a chat
type RenameChatFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ChatInfo represents metadata about a chat session (moved from storage/repository.go)
type ChatInfo struct {
	ID          string
	Name        string
	CreatedAt   int64
	LastUpdated int64
	IsFavorite  bool
}
