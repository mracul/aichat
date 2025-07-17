package main

// modals.go - Contains reusable modal components for the TUI application.
// Centralizes Lipgloss style definitions for consistency and maintainability.

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Centralized modal styles
var (
	modalPromptStyle    = lipgloss.NewStyle().Bold(true)
	modalInputBoxStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1).Width(40)
	modalMessageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	modalBoxStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("203")).Padding(1, 2)
	modalOptionStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(false).Width(8).Align(lipgloss.Center)
	modalOptionSelected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("203")).Background(lipgloss.Color("236"))
	modalInfoTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	modalInfoTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

// ModalType constants for modal type identification.
const (
	ModalTypeInputBox     = "input_box"
	ModalTypeConfirmation = "confirmation"
	ModalTypeInformation  = "information"
)

// InputBoxModal is a reusable modal for text input (e.g. custom chat flow)
type InputBoxModal struct {
	Prompt      string
	Value       string
	Cursor      int
	SelectStart int // -1 if no selection
	SelectEnd   int // -1 if no selection
	Message     string
	Width       int
	Height      int
	Quitting    bool
	closeSelf   func()
}

// Type returns the modal type.
func (m *InputBoxModal) Type() string {
	return ModalTypeInputBox
}

// View renders the input box modal.
func (m *InputBoxModal) View() string {
	input := m.Value
	if m.SelectStart >= 0 && m.SelectEnd > m.SelectStart {
		start, end := m.SelectStart, m.SelectEnd
		if start > end {
			start, end = end, start
		}
		selected := modalInputBoxStyle.Reverse(true).Render(input[start:end])
		input = input[:start] + selected + input[end:]
	}
	if m.Cursor >= 0 && m.Cursor <= len(input) {
		input = input[:m.Cursor] + "|" + input[m.Cursor:]
	}
	prompt := modalPromptStyle.Render(m.Prompt)
	inputBox := modalInputBoxStyle.Render(input)
	msg := ""
	if m.Message != "" {
		msg = "\n" + modalMessageStyle.Render(m.Message)
	}
	content := prompt + "\n" + inputBox + msg
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

// Update handles Bubble Tea messages for the modal.
func (m *InputBoxModal) Update(msg interface{}, ctx interface{}, nav interface{}) (interface{}, interface{}) {
	// Implement message handling as needed, or return self for now
	return m, nil
}

// OnShow is called when the modal is shown.
func (m *InputBoxModal) OnShow() {}

// OnHide is called when the modal is hidden.
func (m *InputBoxModal) OnHide() {}

// IsClosable returns true if the modal can be closed.
func (m *InputBoxModal) IsClosable() bool { return true }

// CloseSelf closes the modal and returns to the previous state.
func (m *InputBoxModal) CloseSelf() {
	if m.closeSelf != nil {
		m.closeSelf()
	}
}

// ConfirmationModal is a reusable yes/no modal for confirmations and notices
type ConfirmationModal struct {
	Title     string
	Prompt    string
	Selected  int // 0 = Yes, 1 = No
	Width     int
	Height    int
	closeSelf func()
}

// Type returns the modal type.
func (m *ConfirmationModal) Type() string {
	return ModalTypeConfirmation
}

// View renders the confirmation modal.
func (m *ConfirmationModal) View() string {
	boxWidth := 40
	prompt := m.Prompt
	options := []string{"Yes", "No"}
	var renderedOptions []string
	for i, opt := range options {
		style := modalOptionStyle
		if i == m.Selected {
			style = modalOptionSelected
		}
		renderedOptions = append(renderedOptions, style.Render(opt))
	}
	optionsLine := lipgloss.JoinHorizontal(lipgloss.Center, renderedOptions...)
	content := lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Center).Render(prompt + "\n\n" + optionsLine)
	box := modalBoxStyle.Width(boxWidth + 4).Align(lipgloss.Center).Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box)
}

// Update handles Bubble Tea messages for the modal.
func (m *ConfirmationModal) Update(msg interface{}, ctx interface{}, nav interface{}) (interface{}, interface{}) {
	// Implement message handling as needed, or return self for now
	return m, nil
}

// OnShow is called when the modal is shown.
func (m *ConfirmationModal) OnShow() {}

// OnHide is called when the modal is hidden.
func (m *ConfirmationModal) OnHide() {}

// IsClosable returns true if the modal can be closed.
func (m *ConfirmationModal) IsClosable() bool { return true }

// CloseSelf closes the modal and returns to the previous state.
func (m *ConfirmationModal) CloseSelf() {
	if m.closeSelf != nil {
		m.closeSelf()
	}
}

// InformationModal is a reusable modal for help/about/info screens
type InformationModal struct {
	Title     string
	Content   string
	Width     int
	Height    int
	Quitting  bool
	closeSelf func()
}

// Type returns the modal type.
func (m *InformationModal) Type() string {
	return ModalTypeInformation
}

// View renders the information modal.
func (m *InformationModal) View() string {
	title := modalInfoTitleStyle.Render(m.Title)
	lines := strings.Split(parseText(m.Content), "\n")
	var renderedLines []string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "|") || strings.Contains(trim, "---") {
			renderedLines = append(renderedLines, modalInfoTextStyle.SetString("").Render(line))
		} else {
			renderedLines = append(renderedLines, modalInfoTextStyle.Render(line))
		}
	}
	content := strings.Join(renderedLines, "\n")
	box := modalBoxStyle.BorderForeground(lipgloss.Color("39")).Width(m.Width - 10).Render(title + "\n" + content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box)
}

// Update handles Bubble Tea messages for the modal.
func (m *InformationModal) Update(msg interface{}, ctx interface{}, nav interface{}) (interface{}, interface{}) {
	// Implement message handling as needed, or return self for now
	return m, nil
}

// OnShow is called when the modal is shown.
func (m *InformationModal) OnShow() {}

// OnHide is called when the modal is hidden.
func (m *InformationModal) OnHide() {}

// IsClosable returns true if the modal can be closed.
func (m *InformationModal) IsClosable() bool { return true }

// CloseSelf closes the modal and returns to the previous state.
func (m *InformationModal) CloseSelf() {
	if m.closeSelf != nil {
		m.closeSelf()
	}
}

// parseText parses markdown and special characters for display in the UI
func parseText(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	re := regexp.MustCompile(`([*_]{1,2}|` + "`" + `)`)
	content = re.ReplaceAllString(content, "")
	reNL := regexp.MustCompile(`\n{3,}`)
	content = reNL.ReplaceAllString(content, "\n\n")
	return content
}

// ErrorModal is a reusable modal for displaying errors
type ErrorModal struct {
	Message  string
	Width    int
	Height   int
	Quitting bool
}

func (m ErrorModal) View() string {
	errorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15")). // white border
		Padding(1, 2).
		Foreground(lipgloss.Color("196")). // red text
		Width(m.Width - 10)
	box := errorStyle.Render(m.Message + "\n\nESC or Enter to close")
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box)
}

func (m ErrorModal) Init() tea.Cmd { return nil }

func (m ErrorModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c", "esc", "enter":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}
