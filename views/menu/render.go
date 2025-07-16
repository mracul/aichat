// render.go: Composes the menu view from AsciiArtView, MenuBoxView, and ControlInfoView.

package menu

import (
	"aichat/types"

	"github.com/charmbracelet/lipgloss"
)

// RenderMenuView composes the menu view from the atomic components.
func RenderMenuView(
	styles interface {
		SelectedStyle() lipgloss.Style
		TextStyle() lipgloss.Style
		DisabledStyle() lipgloss.Style
		ErrorStyle() lipgloss.Style
		HelpStyle() lipgloss.Style
	},
	getMenuTitle func(types.MenuType) string,
	menuView *types.MenuViewState,
	width, height int,
) string {
	if len(types.Menus[menuView.Type].Entries) == 0 {
		return styles.ErrorStyle().Render("Invalid menu type")
	}
	ascii := AsciiArtView()
	menuBox := MenuBoxView(styles, getMenuTitle, menuView, width, height)
	// Pass only HelpStyle to ControlInfoView
	helpStyles := struct{ HelpStyle lipgloss.Style }{HelpStyle: styles.HelpStyle()}
	controlInfo := ControlInfoView(helpStyles, menuView.Type, width)
	titleColor := types.MenuTitleColorMap[menuView.Type]
	titleStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).MarginBottom(1)
	if titleColor != "" {
		titleStyle = titleStyle.Foreground(lipgloss.Color(titleColor))
	} else {
		titleStyle = titleStyle.Foreground(styles.TextStyle().GetForeground())
	}
	title := titleStyle.Render(getMenuTitle(menuView.Type))
	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.PlaceHorizontal(width, lipgloss.Center, ascii),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, title),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, menuBox),
		lipgloss.PlaceHorizontal(width, lipgloss.Left, controlInfo),
	)
}
