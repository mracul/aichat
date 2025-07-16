// box.go: Provides the MenuBoxView component for rendering menu headings and items.

package menu

import (
	"aichat/types"

	"github.com/charmbracelet/lipgloss"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MenuBoxView renders the menu heading and items.
func MenuBoxView(styles interface {
	SelectedStyle() lipgloss.Style
	TextStyle() lipgloss.Style
	DisabledStyle() lipgloss.Style
}, getMenuTitle func(types.MenuType) string, menuView *types.MenuViewState, width, height int) string {
	boxWidth := max(300, width)
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("245")).
		Padding(1, 4).
		Width(boxWidth).
		Height(max(400, height)).
		Align(lipgloss.Center)

	heading := lipgloss.PlaceHorizontal(boxWidth, lipgloss.Center, getMenuTitle(menuView.Type))

	var menuLines []string
	for i, option := range types.Menus[menuView.Type].Entries {
		line := option.Text
		if option.Description != "" {
			line += " - " + option.Description
		}
		var itemStyle lipgloss.Style
		if option.Description != "" {
			itemStyle = lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Left)
		} else {
			itemStyle = lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Center)
		}
		if i == menuView.Selected {
			line = styles.SelectedStyle().Render(line)
		} else {
			line = styles.TextStyle().Render(line)
		}
		if option.Disabled {
			line = styles.DisabledStyle().Render(line)
		}
		menuLines = append(menuLines, itemStyle.Render(line))
	}
	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, heading, "", lipgloss.JoinVertical(lipgloss.Left, menuLines...)))
}
