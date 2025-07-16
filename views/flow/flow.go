// flow.go
// Flow modal system and flow-based view logic for multi-step workflows

package flow

import (
	"aichat/types"

	tea "github.com/charmbracelet/bubbletea"
)

// FlowViewState is a placeholder for a multi-step flow view state
// Implements types.ViewState for integration with navigation/controller

type FlowViewState struct {
	Step int
}

func (f *FlowViewState) Type() string  { return "FlowViewState" }
func (f *FlowViewState) Init() tea.Cmd { return nil }
func (f *FlowViewState) View() string  { return "[Flow Step]" }
func (f *FlowViewState) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// Stub for flow step logic
	return f, nil
}
