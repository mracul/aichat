// control_info.go: Provides the ControlInfoView component for rendering control/help info below menus.

package menu

import (
	"aichat/types"

	"github.com/charmbracelet/lipgloss"
)

func menuMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ControlInfoView renders control info, left-aligned, matching menu box width.
// styles: must have a HelpStyle field (exported)
// menuType: the type of menu for which to render control info
// width: the width to use for rendering the control info
func ControlInfoView(styles interface{ HelpStyle lipgloss.Style }, menuType types.MenuType, width int) string {
	var controlInfo string
	if meta, exists := types.MenuMetas[menuType]; exists {
		if ci, exists := types.ControlInfoMap[meta.ControlInfoType]; exists {
			for _, controlLine := range ci.Lines {
				controlInfo += styles.HelpStyle.Width(menuMax(300, width)).Align(lipgloss.Left).Render(controlLine) + "\n"
			}
		}
	}
	return controlInfo
}
