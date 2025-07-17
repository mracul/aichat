package mvc
package mvc

import (
	"aichat/components/chat/interfaces"
	"aichat/models"

	tea "github.com/charmbracelet/bubbletea"
)

type ChatController struct {
	model models.IChatModel
	view  interfaces.IChatView
}

func NewChatController(model models.IChatModel, view interfaces.IChatView) *ChatController {
	return &ChatController{
		model: model,
		view:  view,
	}
}

// HandleEvent implements interfaces.IChatController
func (c *ChatController) HandleEvent(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement event handling logic, update model as needed
	return c, nil
}

func (c *ChatController) SetView(view interfaces.IChatView) {
	c.view = view
}

func (c *ChatController) SetModel(model models.IChatModel) {
	c.model = model
}

func (c *ChatController) Init() tea.Cmd {
	return nil
}

func (c *ChatController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.HandleEvent(msg)
}

func (c *ChatController) View() string {
	if c.view != nil && c.model != nil {
		return c.view.Render(c.model)
	}
	return "[No view/model set]"
}

