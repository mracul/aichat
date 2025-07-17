package interfaces
package interfaces

type IChatView interface {
	Render(model interface{}) string // Use models.IChatModel in concrete types
}

type IChatController interface {
	HandleEvent(msg interface{}) (interface{}, interface{}) // Use tea.Msg, tea.Model, tea.Cmd in concrete types
	SetView(view IChatView)
	SetModel(model interface{}) // Use models.IChatModel in concrete types
}

