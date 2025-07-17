package chat
package chat

import (
	"aichat/models"

	tea "github.com/charmbracelet/bubbletea"
)

// IChatView defines the rendering interface for chat views.
type IChatView interface {
	Render(model models.IChatModel) string
}

// IChatController defines the event/update logic interface for chat controllers.
type IChatController interface {
	HandleEvent(msg tea.Msg) (tea.Model, tea.Cmd)
	SetView(view IChatView)
	SetModel(model models.IChatModel)
}
