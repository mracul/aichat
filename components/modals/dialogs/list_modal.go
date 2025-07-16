// list_modal.go - Reusable scrollable list modal for Bubble Tea UI
// Displays a windowed list of options (10 at a time), with instruction and control text, styled and centered.

package dialogs

import (
	"aichat/types"
	render "aichat/types/render"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModal is a scrollable, windowed list modal for option selection.
type ListModal struct {
	Title           string
	Options         []string
	Selected        int
	DisplayStart    int    // Index of first displayed entry
	InstructionText string // Shown above the list, left-aligned
	ControlText     string // Shown below the list, left-aligned
	OnSelect        func(index int)
	CloseSelf       func()
	RegionWidth     int // For centering
	RegionHeight    int // For centering
	// [MIGRATION] Use RenderStrategy and Theme for all rendering in ListModal.
	// Replace direct lipgloss.NewStyle() and hardcoded colors with ApplyStrategy and ThemeMap lookups.
	// Add a ThemeMap field to ListModal and use it in ViewRegion().
	ThemeMap render.ThemeMap // Map of styles for rendering
}

// Implement types.ViewState interface
func (m *ListModal) IsMainMenu() bool                   { return false }
func (m *ListModal) MarshalState() ([]byte, error)      { return nil, nil }
func (m *ListModal) UnmarshalState(data []byte) error   { return nil }
func (m *ListModal) ViewType() types.ViewType           { return types.ModalStateType }
func (m *ListModal) GetControlSets() []types.ControlSet { return nil }
func (m *ListModal) Type() types.ViewType               { return types.ModalStateType }

const (
	listWindowSize = 10
	listModalWidth = 40 // Narrower than main menu
)

// Init (Bubble Tea compatibility)
func (m *ListModal) Init() tea.Cmd { return nil }

// Update handles up/down navigation, enter to select, esc to close, and scrolls window as needed.
func (m *ListModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.UpdateWithContext(msg, nil, nil)
}

// UpdateWithContext is a stub for context-aware update logic.
func (m *ListModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// TODO: Implement context-aware update logic
	return m, nil
}

// View renders the modal, centered, with instruction/control text and a windowed list.
func (m *ListModal) View() string {
	return m.ViewRegion(m.RegionWidth, m.RegionHeight)
}

// ViewRegion renders the modal centered in the given region.
func (m *ListModal) ViewRegion(regionWidth, regionHeight int) string {
	// Instruction text (above)
	instr := m.InstructionText
	if instr == "" {
		instr = "Select item:"
	}
	instrLine := lipgloss.NewStyle().MarginBottom(1).Align(lipgloss.Left).Render(instr)

	// List entries (windowed)
	start := m.DisplayStart
	end := start + listWindowSize
	if end > len(m.Options) {
		end = len(m.Options)
	}
	var opts string
	for i := start; i < end; i++ {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.Selected {
			style = style.Bold(true).Foreground(lipgloss.Color("33")).Background(lipgloss.Color("236"))
		}
		opts += style.Render(m.Options[i]) + "\n"
	}
	listBox := lipgloss.NewStyle().Width(listModalWidth).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("245")).Padding(1, 0).Align(lipgloss.Left).Render(opts)

	// Control text (below)
	ctrl := m.ControlText
	if ctrl == "" {
		ctrl = "Up/Down: Navigate, Enter: Select, Esc: Cancel"
	}
	ctrlLine := lipgloss.NewStyle().MarginTop(1).Align(lipgloss.Left).Render(ctrl)

	// Compose modal layout
	content := instrLine + listBox + ctrlLine
	return lipgloss.Place(regionWidth, regionHeight, lipgloss.Center, lipgloss.Center, content)
}
