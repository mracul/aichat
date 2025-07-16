// apikey.go
// Type definitions for API key-related flows (add, remove, set active, test, list)

package flows

import "aichat/types"

// AddAPIKeyFlow for adding a new API key
type AddAPIKeyFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// RemoveAPIKeyFlow for deleting an API key
type RemoveAPIKeyFlow struct {
	ListModal    types.ViewState
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// SetActiveAPIKeyFlow for selecting an active API key
type SetActiveAPIKeyFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// TestAPIKeyFlow for testing the current API key
type TestAPIKeyFlow struct {
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ListAPIKeysFlow for listing and selecting API keys
type ListAPIKeysFlow struct {
	ListModal types.ViewState
	Strategy  FlowStrategy
}
