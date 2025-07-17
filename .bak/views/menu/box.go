package menu
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

	// Calculate max width of menu entries (text + description)
	maxEntryWidth := 0
	entries := types.GetMenuEntries(menuView.MenuType())
	for _, option := range entries {
		entryLen := lipgloss.Width(option.Text)
		if option.Description != "" {
			entryLen += 3 + lipgloss.Width(option.Description) // ' - '
		}
		if entryLen > maxEntryWidth {
			maxEntryWidth = entryLen
		}
	}
	// Center the block of entries in the box
	blockWidth := maxEntryWidth
	if blockWidth > boxWidth-8 { // ensure it doesn't overflow
		blockWidth = boxWidth - 8
	}

	// Heading with larger top margin
	heading := lipgloss.NewStyle().MarginTop(1).Width(boxWidth).Align(lipgloss.Center).Render(getMenuTitle(menuView.MenuType()))

	var menuLines []string
	for i, option := range entries {
		line := option.Text
		// Only show description for main menu
		if menuView.IsMainMenu() && option.Description != "" {
			line += " - " + option.Description
		}
		itemStyle := lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Center)
		if i == menuView.Cursor() {
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
