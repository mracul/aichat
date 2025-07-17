package dialogs
package dialogs

import (
	"aichat/components/modals"
	"aichat/interfaces"

	tea "github.com/charmbracelet/bubbletea"
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
	CloseSelfFunc   func()
	RegionWidth     int                      // For centering
	RegionHeight    int                      // For centering
	Config          modals.ModalRenderConfig // Use ModalRenderConfig for theming/strategy
}

func (m *ListModal) OnShow()          {}
func (m *ListModal) OnHide()          {}
func (m *ListModal) IsClosable() bool { return true }
func (m *ListModal) CloseSelf() {
	if m.CloseSelfFunc != nil {
		m.CloseSelfFunc()
	}
}

const (
	listWindowSize = 10
	listModalWidth = 40 // Narrower than main menu
)

// Init (Bubble Tea compatibility)
func (m *ListModal) Init() tea.Cmd { return nil }

// Update handles up/down navigation, enter to select, esc to close, and scrolls window as needed.
func (m *ListModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
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
			if m.CloseSelfFunc != nil {
				m.CloseSelfFunc()
			}
		case "esc":
			if m.CloseSelfFunc != nil {
				m.CloseSelfFunc()
			}
		}
	}
	return m, nil
}

// UpdateWithContext is a stub for context-aware update logic.
func (m *ListModal) UpdateWithContext(msg tea.Msg, ctx interfaces.Context, nav interfaces.Controller) (tea.Model, tea.Cmd) {
	return m.Update(msg)
}

// View renders the modal, centered, with instruction/control text and a windowed list.
func (m *ListModal) View() string {
	return m.ViewRegion(m.RegionWidth, m.RegionHeight)
}

// ViewRegion renders the modal centered in the given region.
func (m *ListModal) ViewRegion(regionWidth, regionHeight int) string {
	instr := m.InstructionText
	if instr == "" {
		instr = "Select item:"
	}
	// TODO: Integrate proper rendering for modal content here
	// return m.Config.RenderContentWithStrategy(content, "modalBox")
	return "[Modal content rendering not implemented]"
}

// Add ViewState compliance methods
func (m *ListModal) IsMainMenu() bool              { return false }
func (m *ListModal) Type() interfaces.ViewType     { return interfaces.ModalStateType }
func (m *ListModal) ViewType() interfaces.ViewType { return interfaces.ModalStateType }

// Add MarshalState to implement ViewState
func (m *ListModal) MarshalState() ([]byte, error) { return nil, nil }

// Add UnmarshalState to implement ViewState
func (m *ListModal) UnmarshalState(data []byte) error { return nil }

