package models
package models

import (
	"aichat/types"
)

// ChatMessage represents a single message in the chat
// IsUser: true if sent by user, false if assistant
// Content: message text
type ChatMessage struct {
	Content string
	IsUser  bool
}

// IChatModel defines the data access/mutation interface for chat models.
type IChatModel interface {
	GetMessages() []ChatMessage
	AddMessage(msg ChatMessage)
	// Add other data methods as needed
}

// ChatViewState represents the chat window (center pane) in the three-pane layout.
type ChatViewState struct {
	Messages           []ChatMessage
	Streaming          bool
	ChatTitle          string
	ScrollPos          int
	SelectedMessageIdx int
	ResponseReceived   bool
	WaitingForResponse bool
	observers          []types.Observer // Observer pattern
	// ThemeMap, Strategies, MessageList, Spinner, Paginator, Help, ShowHelp are UI-related and should be handled in the view/controller layer.
}

// Ensure ChatViewState implements IChatModel
var _ IChatModel = (*ChatViewState)(nil)

func (c *ChatViewState) GetMessages() []ChatMessage {
	return c.Messages
}

// Observer pattern methods
func (c *ChatViewState) RegisterObserver(o types.Observer) {
	c.observers = append(c.observers, o)
}

func (c *ChatViewState) UnregisterObserver(o types.Observer) {
	for i, obs := range c.observers {
		if obs == o {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

func (c *ChatViewState) NotifyObservers(event interface{}) {
	for _, o := range c.observers {
		o.Notify(event)
	}
}

func (c *ChatViewState) AddMessage(msg ChatMessage) {
	c.Messages = append(c.Messages, msg)
	c.NotifyObservers(types.Event{Type: "message_added", Payload: msg})
}

// chatMessageItem wraps ChatMessage for Bubbles list.Model
type ChatMessageItem struct {
	Msg ChatMessage
}

// NewChatViewState constructs a ChatViewState with a title and messages
func NewChatViewState(title string, messages []ChatMessage) *ChatViewState {
	return &ChatViewState{
		Messages:  messages,
		ChatTitle: title,
	}
}

// ChatMessageItem wraps ChatMessage for Bubbles list.Model
// Implements list.Item interface
func (i ChatMessageItem) FilterValue() string {
	return i.Msg.Content
}

