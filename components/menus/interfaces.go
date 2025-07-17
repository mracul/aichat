package menus
package menus

import "aichat/types"

type IMenuModel interface {
	GetEntries() []types.MenuEntry
	GetCurrentMenuType() types.MenuType
	GetSelectedIndex() int
	SetSelectedIndex(idx int)
}

type IMenuView interface {
	Render(model IMenuModel) string
}

type IMenuController interface {
	HandleEvent(event interface{}) (interface{}, error) // Use concrete event types in implementation
	SetView(view IMenuView)
	SetModel(model IMenuModel)
	MoveSelectionUp()
	MoveSelectionDown()
	Select()
}

