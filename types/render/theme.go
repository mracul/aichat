package render
// theme.go
// Theme type for TUI color schemes, designed for Bubble Tea/lipgloss

package render

import "github.com/charmbracelet/lipgloss"

// DefaultThemeMap provides default styles for modal/dialog rendering
var DefaultThemeMap = map[string]lipgloss.Style{
	"selectedForeground": lipgloss.NewStyle().Foreground(lipgloss.Color("203")),
	"selectedBackground": lipgloss.NewStyle().Background(lipgloss.Color("236")),
	"modalBox":           lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2),
	// Add more as needed
}

// Theme defines a color scheme for a UI component or the whole app
// Colors are hex strings (e.g., "#ffffff")
type Theme struct {
	Name        string // e.g., "Solarized Dark"
	TextColor   string // e.g., "#839496"
	BorderColor string // e.g., "#586e75"
	BgColor     string // optional, e.g., "#002b36"
	AccentColor string // optional, e.g., "#b58900"
}

// ThemeMap maps a theme key (e.g., "menu", "notice") to a Theme
// This allows per-component theming
type ThemeMap map[string]Theme

