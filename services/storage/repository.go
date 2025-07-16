// Package storage provides repository interfaces for persistent data storage.
package storage

import (
	"aichat/types"
	"aichat/types/flows"
	"os"
)

// Extend ChatRepository for file info (modification time)
type ChatRepository interface {
	GetAll() ([]*types.ChatFile, error)
	GetByID(name string) (*types.ChatFile, error)
	Save(chat *types.ChatFile) error
	Delete(name string) error
	GetChatFileInfo(name string) (os.FileInfo, error) // new
}

// PromptRepository defines CRUD operations for prompt templates.
type PromptRepository interface {
	GetAll() ([]*flows.Prompt, error)
	GetByID(name string) (*flows.Prompt, error)
	Save(prompt *flows.Prompt) error
	Delete(name string) error
}

// ModelRepository defines CRUD operations for AI model configs.
type ModelRepository interface {
	GetAll() ([]*types.Model, error)
	GetByID(name string) (*types.Model, error)
	Save(model *types.Model) error
	Delete(name string) error
}
