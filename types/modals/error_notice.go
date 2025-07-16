// error_notice.go - Modal for displaying a centered error message with a single OK confirmation
// Used as a flow or standalone modal after an error in another flow

package modals

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ErrorNoticeModal displays a centered error message and a single OK button
// On pressing OK (Enter), calls OnConfirm and removes itself from view
// Returns no data, just signals completion

type ErrorNoticeModal struct {
	Code      string
	Message   string
	OnConfirm func()
}

func NewErrorNoticeModal(code, message string, onConfirm func()) *ErrorNoticeModal {
	return &ErrorNoticeModal{
		Code:      code,
		Message:   message,
		OnConfirm: onConfirm,
	}
}

func (m *ErrorNoticeModal) Init() tea.Cmd { return nil }

func (m *ErrorNoticeModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEnter {
			if m.OnConfirm != nil {
				m.OnConfirm()
			}
			// Modal manager should remove this modal from view
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ErrorNoticeModal) View() string {
	// Centered error box with code, message, and OK
	return fmt.Sprintf(`
┌───────────────────────────────┐
│                               │
│   ERROR: %s                   │
│   %s                          │
│                               │
│         [ OK ]                │
└───────────────────────────────┘
`, m.Code, m.Message)
}

// Usage:
// modal := NewErrorNoticeModal("401", "Invalid API Key", func() {
//     // OnConfirm: remove modal, return to main menu
//     nav.Replace(mainMenuViewState) // or appropriate navigation action
// })
//
// In your menu/modal manager, ensure OnConfirm navigates to the main menu after error acknowledgment.
