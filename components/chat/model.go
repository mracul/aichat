// model.go - ChatViewState for the center chat window in the three-pane layout.
// Implements ViewState for integration with the navigation stack and app orchestration.
// This state is responsible for rendering chat messages, handling streaming, and managing chat-specific logic.

package chat

import (
	"aichat/types"
	"aichat/types/render"

	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// ViewState interface is imported from types package

// ChatMessage represents a single message in the chat
// IsUser: true if sent by user, false if assistant
// Content: message text
type ChatMessage struct {
	Content string
	IsUser  bool
}

// ChatViewState represents the chat window (center pane) in the three-pane layout.
type ChatViewState struct {
	Messages           []ChatMessage   // Chat messages (now using struct)
	Streaming          bool            // Whether a message is currently streaming
	ChatTitle          string          // Title or summary of the chat
	ScrollPos          int             // Current scroll position (index of first visible message)
	SelectedMessageIdx int             // Index of selected message in visible window, -1 if none
	ResponseReceived   bool            // True if a response was just received
	WaitingForResponse bool            // True if waiting for a response
	ThemeMap           render.ThemeMap // Map of theme names to their styles

	// Bubbles list for messages
	MessageList list.Model

	// Bubbles spinner for streaming
	Spinner spinner.Model

	// Bubbles paginator for message pages
	Paginator paginator.Model

	// Bubbles help for keybindings
	Help help.Model

	// Show help bar when true
	ShowHelp bool
}

// chatMessageItem wraps ChatMessage for Bubbles list.Model
// Implements list.Item interface
type chatMessageItem struct {
	Msg ChatMessage
}

func (i chatMessageItem) Title() string {
	if i.Msg.IsUser {
		return "You"
	}
	return "AI"
}

func (i chatMessageItem) Description() string {
	return i.Msg.Content
}

func (i chatMessageItem) FilterValue() string {
	return i.Msg.Content
}

// chatMessageDelegate is a custom list delegate for chat messages
// Renders markdown for AI messages, highlights selected message
type chatMessageDelegate struct{}

func (d chatMessageDelegate) Height() int                               { return 2 }
func (d chatMessageDelegate) Spacing() int                              { return 1 }
func (d chatMessageDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d chatMessageDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	chatItem, ok := item.(chatMessageItem)
	if !ok {
		return
	}
	var content string
	if !chatItem.Msg.IsUser {
		// AI message: render markdown with ANSI colors using Glamour
		rendered, err := glamour.Render(chatItem.Msg.Content, "dark")
		if err != nil {
			content = chatItem.Msg.Content // fallback to plain text
		} else {
			content = rendered
		}
	} else {
		content = chatItem.Msg.Content
	}
	prefix := "  "
	if index == m.Index() {
		prefix = "> "
	}
	w.Write([]byte(prefix + content))
}

// NewChatViewState constructs a ChatViewState with a Bubbles list.Model, spinner, paginator, and help
func NewChatViewState(title string, messages []ChatMessage) *ChatViewState {
	items := make([]list.Item, len(messages))
	for i, msg := range messages {
		items[i] = chatMessageItem{Msg: msg}
	}
	sp := spinner.New()
	pg := paginator.New()
	pg.PerPage = 10
	hp := help.New()
	delegate := chatMessageDelegate{}
	l := list.New(items, delegate, 60, 20)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	return &ChatViewState{
		Messages:    messages,
		ChatTitle:   title,
		MessageList: l,
		Spinner:     sp,
		Paginator:   pg,
		Help:        hp,
		ShowHelp:    false,
	}
}

// ViewType returns the view type for this state
func (c *ChatViewState) ViewType() types.ViewType {
	return types.ChatStateType
}

// Add Type() method to satisfy ViewState interface
func (c *ChatViewState) Type() types.ViewType {
	return types.ChatStateType
}

// IsMainMenu returns false; this is not the main menu state.
func (c *ChatViewState) IsMainMenu() bool {
	return false
}

// MarshalState serializes the chat view state
func (c *ChatViewState) MarshalState() ([]byte, error) {
	// TODO: Implement full serialization if needed
	return nil, nil
}

// UnmarshalState deserializes the chat view state
func (c *ChatViewState) UnmarshalState(data []byte) error {
	// TODO: Implement full deserialization if needed
	return nil
}

// View renders the chat window using Bubbles list.Model, paginator at the bottom, and help if ShowHelp is true.
func (c *ChatViewState) View() string {
	view := ""
	if c.WaitingForResponse {
		view += c.Spinner.View() + " Waiting for response...\n"
	}
	view += c.MessageList.View()
	if c.Streaming {
		view += "\n[Streaming...]\n"
	}
	if c.ShowHelp {
		view += "\n" + c.Help.View(nil) + "\n"
	}
	view += "\n" + c.Paginator.View() + "\n"
	return view
}

// Update handles Bubble Tea messages, delegating to Bubbles list.Model, spinner, paginator, and help.
func (c *ChatViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement message streaming, chat updates, etc.
	return c, nil
}

// Add UpdateWithContext to satisfy ViewState interface
func (c *ChatViewState) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	model, cmd := c.Update(msg)
	return model.(tea.Model), cmd
}

// Add method to add a message to the list and Messages slice
func (c *ChatViewState) AddMessage(msg ChatMessage) {
	c.Messages = append(c.Messages, msg)
	c.MessageList.InsertItem(len(c.Messages)-1, chatMessageItem{Msg: msg})
}

// GetControlSets returns the chat window's control sets
func (c *ChatViewState) GetControlSets() []types.ControlSet {
	controls := []types.ControlSet{
		{
			Controls: []types.ControlType{
				{
					Name: "Up", Key: tea.KeyUp,
					Action: func() bool {
						if c.ScrollPos > 0 {
							c.ScrollPos--
							return true
						}
						return false
					},
				},
				{
					Name: "Down", Key: tea.KeyDown,
					Action: func() bool {
						c.ScrollPos++ // TODO: bound by message count
						return true
					},
				},
				{
					Name: "PageUp", Key: tea.KeyPgUp,
					Action: func() bool {
						c.ScrollPos -= 10
						if c.ScrollPos < 0 {
							c.ScrollPos = 0
						}
						return true
					},
				},
				{
					Name: "PageDown", Key: tea.KeyPgDown,
					Action: func() bool {
						c.ScrollPos += 10 // TODO: bound by message count
						return true
					},
				},
				{
					Name: "Enter", Key: tea.KeyEnter,
					Action: func() bool {
						// TODO: select message
						return true
					},
				},
				{
					Name: "Esc", Key: tea.KeyEsc,
					Action: func() bool {
						// TODO: unfocus logic
						return true
					},
				},
			},
		},
	}
	return controls
}

// GetControlSet returns the chat window's control set (legacy method)
func (c *ChatViewState) GetControlSet() interface{} {
	return nil // No longer using a legacy control set
}

// Implement Init() method for ChatViewState to satisfy ViewState interface
func (c *ChatViewState) Init() tea.Cmd {
	return nil
}

// TODO: Integrate with sidebar navigation and input area for full three-pane layout.
// TODO: Support message streaming, scrolling, and chat metadata display.
// TODO: Replace string messages with proper message structs.
