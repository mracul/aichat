package menu
// render.go: Composes the menu view from AsciiArtView, MenuBoxView, and ControlInfoView.

package menu

import (
	"aichat/types"

	"github.com/charmbracelet/lipgloss"
)

type helpStyleWrapper struct {
	s interface{ HelpStyle() lipgloss.Style }
}

func (h helpStyleWrapper) HelpStyle() lipgloss.Style { return h.s.HelpStyle() }

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
	if len(types.GetMenuEntries(menuView.MenuType())) == 0 {
		return styles.ErrorStyle().Render("Invalid menu type")
	}
	ascii := AsciiArtView()
	menuBox := MenuBoxView(styles, getMenuTitle, menuView, width, height)
	controlInfo := ControlInfoView(helpStyleWrapper{styles}, menuView.MenuType(), width)
	titleColor := types.MenuTitleColorMap[menuView.MenuType()]
	titleStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).MarginBottom(1)
	if titleColor != "" {
		titleStyle = titleStyle.Foreground(lipgloss.Color(titleColor))
	} else {
		titleStyle = titleStyle.Foreground(styles.TextStyle().GetForeground())
	}
	title := titleStyle.Render(getMenuTitle(menuView.MenuType()))
	// Instead of centering, align control info with menu box left edge
	boxWidth := width
	boxLeft := 0
	if menuBoxWidth := 0; len(menuBox) > 0 {
		menuBoxWidth = lipgloss.Width(menuBox)
		if menuBoxWidth > 0 {
			boxWidth = menuBoxWidth
			boxLeft = (width - boxWidth) / 2
		}
	}
	controlInfoStyled := lipgloss.NewStyle().MarginTop(1).Width(boxWidth).Align(lipgloss.Left).Render(controlInfo)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.PlaceHorizontal(width, lipgloss.Center, ascii),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, title),
		lipgloss.PlaceHorizontal(width, lipgloss.Center, menuBox),
		lipgloss.PlaceHorizontal(width, lipgloss.Left, lipgloss.NewStyle().MarginLeft(boxLeft).Render(controlInfoStyled)),
	)
}

