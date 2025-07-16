// prompts.go
// Type definitions for prompts-related flows (add, remove, set default, list)

package flows

import (
	"aichat/types"
	"encoding/json"
	"os"
)

// Prompt represents a prompt template for the assistant.
type Prompt struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Default bool   `json:"default,omitempty"`
}

// AddPromptFlow for creating a new prompt
type AddPromptFlow struct {
	InputModal  types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// RemovePromptFlow for deleting a prompt
type RemovePromptFlow struct {
	ListModal    types.ViewState
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// SetDefaultPromptFlow for selecting a default prompt
type SetDefaultPromptFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}

// ListPromptsFlow for listing and selecting prompts
type ListPromptsFlow struct {
	ListModal types.ViewState
	Strategy  FlowStrategy
}

// SavePromptsToFile saves a slice of prompts to a JSON file
func SavePromptsToFile(prompts []Prompt, filePath string) error {
	data, err := json.MarshalIndent(prompts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
