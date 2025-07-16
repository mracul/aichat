// menu.go - Contains MenuModal for displaying a menu with selectable options in a modal dialog in the Bubble Tea UI.

package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuModal is a reusable modal for displaying a menu with options.
type MenuModal struct {
	Title        string
	Options      []string
	Selected     int
	OnSelect     func(index int)
	CloseSelf    func()
	RegionWidth  int // Last-known or intended region width for rendering
	RegionHeight int // Last-known or intended region height for rendering
	// [MIGRATION] Use RenderStrategy and Theme for all rendering in MenuModal.
	// Replace direct lipgloss.NewStyle() and hardcoded colors with ApplyStrategy and ThemeMap lookups.
	// Add a ThemeMap field to MenuModal and use it in ViewRegion().
	ThemeMap       map[string]lipgloss.Style
	RenderStrategy func(string, int, int) string
}

// Init initializes the modal (Bubble Tea compatibility).
func (m *MenuModal) Init() tea.Cmd { return nil }

// Update handles Bubble Tea messages for the modal: up/down to navigate, enter to select, esc to close.
func (m *MenuModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Selected > 0 {
				m.Selected--
			} else {
				m.Selected = len(m.Options) - 1
			}
		case "down":
			if m.Selected < len(m.Options)-1 {
				m.Selected++
			} else {
				m.Selected = 0
			}
		case "enter":
			if m.OnSelect != nil {
				m.OnSelect(m.Selected)
			}
			if m.CloseSelf != nil {
				m.CloseSelf()
			}
			return m, tea.Quit
		case "esc":
			if m.CloseSelf != nil {
				m.CloseSelf()
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the menu modal UI as a string, centered in the stored region (RegionWidth, RegionHeight).
// The title is shown above, and options are rendered vertically, with the selected option highlighted.
func (m *MenuModal) View() string {
	return m.ViewRegion(m.RegionWidth, m.RegionHeight)
}

// ViewRegion renders the menu modal UI as a string, centered in the given region (width, height).
func (m *MenuModal) ViewRegion(regionWidth, regionHeight int) string {
	title := m.ThemeMap["title"].Render(m.Title)
	var opts string
	for i, opt := range m.Options {
		style := m.ThemeMap["option"]
		if i == m.Selected {
			// Instead of Foreground/Background expecting TerminalColor, use a color constant or skip
			style = style.Bold(true)
		}
		opts += style.Render(opt) + "\n"
	}
	content := title + "\n\n" + opts
	return m.RenderStrategy(content, regionWidth, regionHeight)
}
