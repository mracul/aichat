package chatwindow

import (
	"aichat/types"
	"encoding/json"

	render "aichat/types/render"

	tea "github.com/charmbracelet/bubbletea"
)

// ChatWindowViewState represents the state of the chat window (messages, input, etc.)
type ChatWindowViewState struct {
	ChatID      string
	Messages    []types.Message
	Metadata    types.ChatMetadata
	InputBuffer string
	Focus       string // "chat", "input", etc.
	// [MIGRATION] Use RenderStrategy and Theme for all rendering in ChatWindowViewState.
	// Replace direct rendering logic with ApplyStrategy and ThemeMap lookups.
	// Add a ThemeMap field to ChatWindowViewState and use it in ViewMessages() and ViewInput().
	ThemeMap render.ThemeMap
}

func (c *ChatWindowViewState) Update(msg tea.Msg) (types.ViewState, tea.Cmd) {
	// TODO: Implement chat message/input event handling
	return c, nil
}

func (c *ChatWindowViewState) View() string {
	return c.ViewMessages() + "\n" + c.ViewInput()
}

func (c *ChatWindowViewState) MarshalState() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ChatWindowViewState) UnmarshalState(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *ChatWindowViewState) ViewType() types.ViewType { return types.ChatStateType }
func (c *ChatWindowViewState) IsMainMenu() bool         { return false }

// Add Type() method to satisfy ViewState interface
func (c *ChatWindowViewState) Type() types.ViewType {
	return types.ChatStateType
}

// Add UpdateWithContext to satisfy ViewState interface
func (c *ChatWindowViewState) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	model, cmd := c.Update(msg)
	return model.(tea.Model), cmd
}

// GetControlSets returns the chat window's control sets
func (c *ChatWindowViewState) GetControlSets() []types.ControlSet {
	controls := []types.ControlSet{
		{
			Controls: []types.ControlType{
				{
					Name: "Tab", Key: tea.KeyTab,
					Action: func() bool {
						// TODO: switch focus between chat and input
						return true
					},
				},
				{
					Name: "Esc", Key: tea.KeyEsc,
					Action: func() bool {
						// TODO: handle back/cancel
						return true
					},
				},
			},
		},
	}
	return controls
}

// ViewMessages renders the chat messages area.
func (c *ChatWindowViewState) ViewMessages() string {
	// TODO: Render chat messages
	return "[Chat messages here]"
}

// ViewInput renders the input area.
func (c *ChatWindowViewState) ViewInput() string {
	// TODO: Render input area
	return "[Input area here]"
}

// Implement Init() method for ChatWindowViewState
func (c *ChatWindowViewState) Init() tea.Cmd {
	return nil
}
