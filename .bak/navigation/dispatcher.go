package navigation
package navigation

import (
	"aichat/types"

	tea "github.com/charmbracelet/bubbletea"
)

// NavigationController centralizes all navigation actions for menus, modals, and view states.
type NavigationController interface {
	// NavigateTo pushes a new menu onto the stack, setting the current menu as parent.
	NavigateTo(next types.MenuType, from types.ViewState)
	// Pop removes the current view from the stack, returning to the previous one.
	Pop()
	// ReplaceCurrent replaces the current view with a new one.
	ReplaceCurrent(next types.ViewState)
	// ShowModal displays a modal dialog (confirmation, input, etc.).
	ShowModal(modalType types.ModalType, data interface{})
	// CurrentState returns the current top view state.
	CurrentState() types.ViewState
}

// AppModel is the root Bubble Tea model with navigation stack.
type AppModel struct {
	Stack *NavigationStack
	// ... other fields ...
}

// Init implements tea.Model's Init method.
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update handles navigation messages and delegates to the top view.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch nav := msg.(type) {
	case NavigationMsg:
		switch nav.Action {
		case PushAction:
			m.Stack.Push(nav.Target)
		case PopAction:
			m.Stack.Pop()
		case ResetAction:
			// TODO: Reset to main menu
		}
		return m, nil
	}
	// Delegate to top view
	top := m.Stack.Top()
	newState, cmd := top.Update(msg)
	// Ensure newState is a types.ViewState
	if vs, ok := newState.(types.ViewState); ok {
		m.Stack.ReplaceTop(vs)
	} else {
		m.Stack.ReplaceTop(top)
	}
	return m, cmd
}

// View renders the current top view.
func (m *AppModel) View() string {
	return m.Stack.Top().View()
}
