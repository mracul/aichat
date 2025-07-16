package types

import (
	"encoding/json"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ModalType is a type alias for modal identifiers (e.g., confirmation, input, etc.)
type ModalType string

// Message represents a chat message (shared across app, for JSON serialization).
type Message struct {
	Role          string `json:"role"`
	Content       string `json:"content"`
	MessageNumber int    `json:"message_number"`
}

// ChatMetadata stores additional information about a chat session.
type ChatMetadata struct {
	Summary    string    `json:"summary,omitempty"`
	Title      string    `json:"title,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Model      string    `json:"model,omitempty"`
	Favorite   bool      `json:"favorite,omitempty"`
	ModifiedAt int64     `json:"modified_at,omitempty"` // Unix timestamp for last modification
}

// ChatFile represents the complete chat file structure for JSON storage.
type ChatFile struct {
	Metadata ChatMetadata `json:"metadata"`
	Messages []Message    `json:"messages"`
}

// Model represents an AI model configuration.
type Model struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// ModelsConfig represents the models configuration stored in JSON.
type ModelsConfig struct {
	Models []Model `json:"models"`
}

// Control represents a generalized control input (key binding, description, and action).
type Control struct {
	Key         string      // The input key (e.g., "ctrl+c", "enter")
	Description string      // Human-readable description of the control
	Action      func() bool // The action to execute (returns true if handled)
}

// ControlAction is a function signature for control actions.
type ControlAction func() bool

// ControlInfo represents the outline of controls for a menu (displayed below the menu box)
type ControlInfo struct {
	Lines []string // Each line is a row of control info (e.g., "↑↓ navigate, Enter select, Esc back")
}

// Predefined control info for menus
var DefaultControlInfo = ControlInfo{
	Lines: []string{"↑↓ navigate", "Enter select", "Esc back"},
}

var FavoritesControlInfo = ControlInfo{
	Lines: []string{"↑↓ navigate", "Enter load chat", "F unfavorite", "Esc back"},
}

// ControlInfoType enumerates the types of control info layouts
// Used to select the correct control info for each menu or entry
type ControlInfoType int

const (
	DefaultControlInfoType ControlInfoType = iota
	FavoritesControlInfoType
	ListChatsControlInfoType
)

// ControlInfoMap maps ControlInfoType to the actual ControlInfo
var ControlInfoMap = map[ControlInfoType]ControlInfo{
	DefaultControlInfoType:   {Lines: []string{"↑↓ navigate", "Enter select", "Esc back"}},
	FavoritesControlInfoType: {Lines: []string{"↑↓ navigate", "Enter load chat", "F unfavorite", "Esc back"}},
	ListChatsControlInfoType: {Lines: []string{"↑↓ navigate", "Enter view chat", "Esc back"}},
}

// MenuMeta holds metadata for each menu, including its control info type
// This allows each menu to display the correct control outline
type MenuMeta struct {
	ControlInfoType ControlInfoType
}

// MenuMetas maps each MenuType to its metadata (including control info type)
// By default, menus use DefaultControlInfoType unless specified otherwise
var MenuMetas = map[MenuType]MenuMeta{
	MainMenu:    {ControlInfoType: DefaultControlInfoType},
	ChatsMenu:   {ControlInfoType: DefaultControlInfoType},
	PromptsMenu: {ControlInfoType: DefaultControlInfoType},
	ModelsMenu:  {ControlInfoType: DefaultControlInfoType},
	APIKeyMenu:  {ControlInfoType: DefaultControlInfoType},
	HelpMenu:    {ControlInfoType: DefaultControlInfoType},
	ExitMenu:    {ControlInfoType: DefaultControlInfoType},
}

// MenuTitleColorMap maps each MenuType to a color string for the menu title.
// This allows menu titles to be themed per menu type.
var MenuTitleColorMap = map[MenuType]string{
	MainMenu:    "33",  // Cyan/blue for main menu
	ChatsMenu:   "129", // Purple for chats
	PromptsMenu: "36",  // Teal for prompts
	ModelsMenu:  "39",  // Blue for models
	APIKeyMenu:  "208", // Orange for API keys
	HelpMenu:    "244", // Gray for help
	ExitMenu:    "196", // Red for exit
}

// ControlType represents a single control action (e.g., Up, Down, Enter, Esc) with an associated action
// Name: human-readable name, Key: key binding, Action: function to execute
// Action can be nil or a function with a standard signature (e.g., func() bool)
type ControlType struct {
	Name   string
	Key    tea.KeyType
	Action func() bool // or interface{} for more flexibility
}

// ControlSet represents a set of controls for a view/state
// Each ControlSet can be customized per view/state
// Example: {ControlUp, ControlDown, ControlEnter, ControlEsc}
type ControlSet struct {
	Controls []ControlType
}

// Example: DefaultControlSet for menus
var DefaultControlSet = ControlSet{
	Controls: []ControlType{
		{Name: "Up", Key: tea.KeyUp, Action: nil},
		{Name: "Down", Key: tea.KeyDown, Action: nil},
		{Name: "Enter", Key: tea.KeyEnter, Action: nil},
		{Name: "Esc", Key: tea.KeyEsc, Action: nil},
	},
}

// APIKey represents a single API key with a title, key, URL, and active status.
type APIKey struct {
	Title  string `json:"title"`
	Key    string `json:"key"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

// APIKeysConfig represents the configuration for multiple API keys.
type APIKeysConfig struct {
	Keys []APIKey `json:"keys"`
}

// ErrorResponse represents an error response from an API.
type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
}

// StreamRequestBody represents the request body for streaming chat completions.
type StreamRequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// SaveAPIKeysToFile saves API keys configuration to a JSON file.
func SaveAPIKeysToFile(config APIKeysConfig, filePath string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

// SaveModelsToFile saves a slice of models to a JSON file.
func SaveModelsToFile(models []Model, filePath string) error {
	config := ModelsConfig{Models: models}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
