package menus
package menus

import (
	"fmt"
)

type MenuView struct{}

func NewMenuView() *MenuView { return &MenuView{} }

func (v *MenuView) Render(model IMenuModel) string {
	entries := model.GetEntries()
	selected := model.GetSelectedIndex()
	out := ""
	for i, entry := range entries {
		prefix := "  "
		if i == selected {
			prefix = "> "
		}
		out += fmt.Sprintf("%s%s - %s\n", prefix, entry.Text, entry.Description)
	}
	return out
}

