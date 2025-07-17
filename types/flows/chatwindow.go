package flows
// chatwindow.go
// Type definitions for chat window flows (export, rename, search, threading)

package flows

import "aichat/types"

// SearchChatFlow for searching within a chat
type SearchChatFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ThreadingFlow for managing message threads (future expansion)
type ThreadingFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

