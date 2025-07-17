package dialogs
// menu.go - Contains MenuModal for displaying a menu with selectable options in a modal dialog in the Bubble Tea UI.

package dialogs

import (
	"aichat/components/modals"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuModal is a reusable modal for displaying a menu with options.
type MenuModal struct {
	Title        string
	Options      []string
	Selected     int
	OnSelect     func(index int)
	CloseSelf    func()
	RegionWidth  int                      // Last-known or intended region width for rendering
	RegionHeight int                      // Last-known or intended region height for rendering
	Config       modals.ModalRenderConfig // Use ModalRenderConfig for theming/strategy
	modals.BaseModal
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
	var opts string
	for i, opt := range m.Options {
		key := "option"
		if i == m.Selected {
			key = "optionSelected"
		}
		opts += m.Config.RenderContentWithStrategy(opt, key) + "\n"
	}
	content := m.Title + "\n\n" + opts
	return m.Config.RenderContentWithStrategy(content, "modalBox")
}

