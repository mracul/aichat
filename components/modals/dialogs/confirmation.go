// confirmation.go - Contains the ConfirmationModal for displaying confirmation dialogs with 1-3 options in the Bubble Tea UI.
// Update logic supports left/right navigation, enter to select, esc to close/cancel.

package dialogs

import (
	"aichat/components/modals"
	"aichat/types"
	render "aichat/types/render"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationModal is a reusable modal for confirmation dialogs (1-3 options).
type ConfirmationModal struct {
	modals.BaseModal
	RegionWidth  int
	RegionHeight int
	focused      bool
	// [MIGRATION] Use RenderStrategy and Theme for all rendering in ConfirmationModal.
	// Replace direct lipgloss.NewStyle() and hardcoded colors with ApplyStrategy and ThemeMap lookups.
	// Add a ThemeMap field to ConfirmationModal and use it in ViewRegion().
	ThemeMap map[string]lipgloss.Style
}

// NewConfirmationModal creates a new ConfirmationModal with the given message, options, and closeSelf callback.
func NewConfirmationModal(message string, options []modals.ModalOption, closeSelf modals.CloseSelfFunc) *ConfirmationModal {
	if len(options) < 1 || len(options) > 3 {
		panic("ConfirmationModal must have 1-3 options")
	}
	return &ConfirmationModal{
		BaseModal: modals.BaseModal{
			Message:   message,
			Options:   options,
			CloseSelf: closeSelf,
			Selected:  0,
		},
		RegionWidth:  60,
		RegionHeight: 10,
		focused:      true,
		ThemeMap:     render.DefaultThemeMap, // Initialize ThemeMap
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
	// [MIGRATION] Use RenderStrategy and Theme for all rendering in ConfirmationModal.
	// Replace direct lipgloss.NewStyle() and hardcoded colors with ApplyStrategy and ThemeMap lookups.
	// Add a ThemeMap field to ConfirmationModal and use it in ViewRegion().
	modalBox := m.ThemeMap["modalBox"].Width(regionWidth).Height(regionHeight)
	msg := lipgloss.NewStyle().Bold(true).Render(m.Message)
	var opts string
	for i, opt := range m.Options {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.Selected {
			style = style.Bold(true).Foreground(lipgloss.Color("33")).Background(lipgloss.Color("236"))
		}
		opts += style.Render(opt.Label)
	}
	content := msg + "\n\n" + opts
	box := modalBox.Render(content)
	return lipgloss.Place(regionWidth, regionHeight, lipgloss.Center, lipgloss.Center, box)
}

// --- ViewState interface ---
func (m *ConfirmationModal) GetControlSets() []types.ControlSet {
	return []types.ControlSet{
		{
			Controls: []types.ControlType{
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
func (m *ConfirmationModal) ViewType() types.ViewType         { return types.ModalStateType }
func (m *ConfirmationModal) Type() types.ViewType             { return types.ModalStateType }

func (m *ConfirmationModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// TODO: Implement context-aware update logic
	return m, nil
}
