package flows
// global.go
// Type definitions for global flows (confirmation, input, selection, error/info)

package flows

import "aichat/types"

// ConfirmationFlow for confirming destructive actions
type ConfirmationFlow struct {
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// InputFlow for collecting user input (name, prompt, etc.)
type InputFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// SelectionFlow for choosing from a list
type SelectionFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ErrorInfoFlow for displaying errors or info messages
type ErrorInfoFlow struct {
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

