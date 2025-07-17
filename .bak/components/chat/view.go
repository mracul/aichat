package chat
package chat

import (
	"aichat/models"
	"strings"
)

type ChatView struct{}

func NewChatView() *ChatView {
	return &ChatView{}
}

// Render renders the chat messages as a string.
func (v *ChatView) Render(model models.IChatModel) string {
	messages := model.GetMessages()
	var sb strings.Builder
	sb.WriteString("[Chat]\n")
	for _, msg := range messages {
		if msg.IsUser {
			sb.WriteString("You: ")
		} else {
			sb.WriteString("AI: ")
		}
		sb.WriteString(msg.Content)
		sb.WriteString("\n")
	}
	return sb.String()
}
