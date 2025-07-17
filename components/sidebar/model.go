package sidebar
// components/sidebar/model.go - SidebarTabsModel using vendored soft-serve tabs implementation

package sidebar

import (
	"aichat/components/sidebar/tabs"

	tea "github.com/charmbracelet/bubbletea"
)

// SidebarTabsModel manages the horizontal chat tabs using the new Tabs component
// Implements tea.Model

type SidebarTabsModel struct {
	Tabs *tabs.Tabs
}

// NewSidebarTabsModel creates a new sidebar tabs model with the given chat names
func NewSidebarTabsModel(tabNames []string) *SidebarTabsModel {
	t := tabs.New(tabNames)
	return &SidebarTabsModel{Tabs: t}
}

// Init implements tea.Model
func (m *SidebarTabsModel) Init() tea.Cmd {
	return m.Tabs.Init()
}

// Update implements tea.Model
func (m *SidebarTabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.Tabs.Update(msg)
	m.Tabs = model.(*tabs.Tabs)
	return m, cmd
}

// View implements tea.Model
func (m *SidebarTabsModel) View() string {
	return m.Tabs.View()
}

// ActiveTab returns the index of the currently active tab
func (m *SidebarTabsModel) ActiveTab() int {
	return m.Tabs.ActiveTab()
}

// TabNames returns the names of all tabs
func (m *SidebarTabsModel) TabNames() []string {
	return m.Tabs.TabNames()
}

