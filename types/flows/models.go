// models.go
// Type definitions for models-related flows (add, remove, set default, list)

package flows

import "aichat/types"

// AddModelFlow for creating a new model
type AddModelFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// RemoveModelFlow for deleting a model
type RemoveModelFlow struct {
	ListModal    types.ViewState
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// SetDefaultModelFlow for selecting a default model
type SetDefaultModelFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ListModelsFlow for listing and selecting models
type ListModelsFlow struct {
	ListModal types.ViewState
	Strategy  FlowStrategy
}
