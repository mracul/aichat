package dialogs
// confirmation.go - Contains the ConfirmationModal for displaying confirmation dialogs with 1-3 options in the Bubble Tea UI.
// Update logic supports left/right navigation, enter to select, esc to close/cancel.

package dialogs

import (
	"aichat/components/modals"
	"aichat/interfaces"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmationModal is a reusable modal for confirmation dialogs (1-3 options).
type ConfirmationModal struct {
	modals.BaseModal
	focused bool
}

func NewConfirmationModal(message string, options []modals.ModalOption, closeSelf modals.CloseSelfFunc, config modals.ModalRenderConfig) *ConfirmationModal {
	if len(options) < 1 || len(options) > 3 {
		panic("ConfirmationModal must have 1-3 options")
	}
	return &ConfirmationModal{
		BaseModal: modals.BaseModal{
			ModalRenderConfig: config,
			Message:           message,
			Options:           options,
			CloseSelf:         closeSelf,
			Selected:          0,
			RegionWidth:       300,
			RegionHeight:      150,
		},
		focused: true,
	}
}

func (m *ConfirmationModal) Init() tea.Cmd { return nil }

func (m *ConfirmationModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.UpdateWithContext(msg, nil, nil)
}

func (m *ConfirmationModal) View() string {
	return m.ViewRegion(m.RegionWidth, m.RegionHeight)
}

func (m *ConfirmationModal) ViewRegion(regionWidth, regionHeight int) string {
	header := m.Message
	// Center the message

	// Render options horizontally, centered, with spacing
	var options []string
	for i, opt := range m.Options {
		label := opt.Label
		if i == m.Selected {
			// Highlight selected option with brackets and theme highlight
			label = "[" + label + "]"
			label = m.RenderContentWithStrategy(label, "highlight")
		} else {
			label = m.RenderContentWithStrategy(label, "modalBox")
		}
		options = append(options, label)
	}
	// Join options with spacing
	optionsLine := "    " + options[0] + "         " + options[1]

	// Compose content: message, blank line, options
	content := "\n" + header + "\n\n" + optionsLine + "\n"
	// Render the modal box with the composed content
	return m.RenderContentWithStrategy(content, "modalBox")
}

// --- ViewState interface ---
func (m *ConfirmationModal) GetControlSets() []interfaces.ControlSet {
	return []interfaces.ControlSet{
		{
			Controls: []interfaces.ControlType{
				{Name: "Left", Key: tea.KeyLeft, Action: func() bool {
					m.BaseModal.Selected = (m.BaseModal.Selected + len(m.BaseModal.Options) - 1) % len(m.BaseModal.Options)
					return true
				}},
				{Name: "Right", Key: tea.KeyRight, Action: func() bool { m.BaseModal.Selected = (m.BaseModal.Selected + 1) % len(m.BaseModal.Options); return true }},
				{Name: "Enter", Key: tea.KeyEnter, Action: func() bool {
					if m.BaseModal.Selected >= 0 && m.BaseModal.Selected < len(m.BaseModal.Options) {
						m.BaseModal.Options[m.BaseModal.Selected].OnSelect()
						if m.BaseModal.CloseSelf != nil {
							m.BaseModal.CloseSelf()
						}
						return true
					}
					return false
				}},
				{Name: "Esc", Key: tea.KeyEsc, Action: func() bool {
					if m.BaseModal.CloseSelf != nil {
						m.BaseModal.CloseSelf()
						return true
					}
					return false
				}},
			},
		},
	}
}
func (m *ConfirmationModal) IsMainMenu() bool                 { return false }
func (m *ConfirmationModal) MarshalState() ([]byte, error)    { return nil, nil }
func (m *ConfirmationModal) UnmarshalState(data []byte) error { return nil }
func (m *ConfirmationModal) ViewType() interfaces.ViewType    { return interfaces.ModalStateType }
func (m *ConfirmationModal) Type() interfaces.ViewType        { return interfaces.ModalStateType }

func (m *ConfirmationModal) UpdateWithContext(msg tea.Msg, ctx interfaces.Context, nav interfaces.Controller) (tea.Model, tea.Cmd) {
	// TODO: Implement context-aware update logic
	return m, nil
}

