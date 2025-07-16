// strategy.go
// Decorator and utility for applying render strategies to views (Bubble Tea/lipgloss)

package render

import "github.com/charmbracelet/lipgloss"

// Renderable is an interface for any view that can be rendered with a strategy
//
type Renderable interface {
	RenderWithStrategy(strategy RenderStrategy, theme Theme) string
}

// ApplyStrategy applies a RenderStrategy and Theme to a string (the view content)
// Returns the styled string for display
func ApplyStrategy(content string, strategy RenderStrategy, theme Theme) string {
	style := lipgloss.NewStyle()
	if strategy.Dimension.Width > 0 {
		style = style.Width(strategy.Dimension.Width)
	}
	if strategy.Dimension.Height > 0 {
		style = style.Height(strategy.Dimension.Height)
	}
	if strategy.Padding > 0 {
		style = style.Padding(strategy.Padding)
	}
	if strategy.Margin > 0 {
		style = style.Margin(strategy.Margin)
	}
	if strategy.Border {
		style = style.Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.BorderColor))
	}
	if strategy.Centered {
		style = style.Align(lipgloss.Center)
	}
	if theme.TextColor != "" {
		style = style.Foreground(lipgloss.Color(theme.TextColor))
	}
	if theme.BgColor != "" {
		style = style.Background(lipgloss.Color(theme.BgColor))
	}
	return style.Render(content)
}
