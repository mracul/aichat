package chatwindow
package chatwindow

import (
	"aichat/components/input"
	"aichat/models"
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
	ThemeMap       render.ThemeMap
	RenderStrategy render.RenderStrategy
	// Add InputModel field to ChatWindowViewState
	InputModel *models.InputModel
}

func (c *ChatWindowViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	// Use centralized theme and strategy for chat messages
	chatTheme, ok := c.ThemeMap["chat"]
	if !ok {
		chatTheme = render.Theme{Name: "default", TextColor: "#ffffff", BgColor: "#000000"}
	}
	return render.ApplyStrategy("[Chat messages here]", c.RenderStrategy, chatTheme)
}

// ViewInput renders the input area.
func (c *ChatWindowViewState) ViewInput() string {
	// Use centralized theme and strategy for input area
	if c.InputModel != nil {
		return input.RenderInputView(c.InputModel, nil) // Pass real context if available
	}
	inputTheme, ok := c.ThemeMap["input"]
	if !ok {
		inputTheme = render.Theme{Name: "default", TextColor: "#ffffff", BgColor: "#222222"}
	}
	return render.ApplyStrategy("[Input area here]", c.RenderStrategy, inputTheme)
}

// Implement Init() method for ChatWindowViewState
func (c *ChatWindowViewState) Init() tea.Cmd {
	return nil
}

// Observer pattern: implement types.Observer
func (c *ChatWindowViewState) Notify(event interface{}) {
	if ev, ok := event.(types.Event); ok {
		switch ev.Type {
		case "message_added":
			// Optionally update Messages or trigger re-render
		case "input_changed":
			// Optionally update InputBuffer or trigger re-render
		}
	}
}

// Register as observer in initialization (example)
func NewChatWindowViewState(chatID string, messages []types.Message, metadata types.ChatMetadata, inputModel *models.InputModel, themeMap render.ThemeMap, strategy render.RenderStrategy) *ChatWindowViewState {
	c := &ChatWindowViewState{
		ChatID:         chatID,
		Messages:       messages,
		Metadata:       metadata,
		InputModel:     inputModel,
		ThemeMap:       themeMap,
		RenderStrategy: strategy,
	}
	// Register as observer to ChatViewState and InputModel if available
	// (Assume you have access to those models here)
	// chatViewState.RegisterObserver(c)
	if inputModel != nil {
		inputModel.RegisterObserver(c)
	}
	return c
}
