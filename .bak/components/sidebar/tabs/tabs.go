package tabs
package tabs

import (
	"strings"

	"aichat/types/render"

	tea "github.com/charmbracelet/bubbletea"
)

// SelectTabMsg is a message that contains the index of the tab to select.
type SelectTabMsg int

// ActiveTabMsg is a message that contains the index of the current active tab.
type ActiveTabMsg int

// Tabs is a Bubble Tea component that displays a list of tabs.
type Tabs struct {
	tabs         []string
	activeTab    int
	TabSeparator string
	UseDot       bool
	Width        int
	ThemeMap     render.ThemeMap
	Strategy     render.RenderStrategy
}

// New creates a new Tabs component with default styles.
func New(tabNames []string) *Tabs {
	return &Tabs{
		tabs:         tabNames,
		activeTab:    0,
		TabSeparator: " | ",
		UseDot:       false,
		Width:        0,
		ThemeMap:     nil,                     // Should be set after construction
		Strategy:     render.RenderStrategy{}, // Should be set after construction
	}
}

// SetSize sets the width for rendering.
func (t *Tabs) SetSize(width int) {
	t.Width = width
}

// Init implements tea.Model.
func (t *Tabs) Init() tea.Cmd {
	t.activeTab = 0
	return nil
}

// Update implements tea.Model.
func (t *Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			t.activeTab = (t.activeTab + 1) % len(t.tabs)
			cmds = append(cmds, t.activeTabCmd)
		case "shift+tab":
			t.activeTab = (t.activeTab - 1 + len(t.tabs)) % len(t.tabs)
			cmds = append(cmds, t.activeTabCmd)
		}
	case SelectTabMsg:
		tab := int(msg)
		if tab >= 0 && tab < len(t.tabs) {
			t.activeTab = tab
		}
	}
	return t, tea.Batch(cmds...)
}

// View implements tea.Model.
func (t *Tabs) View() string {
	s := strings.Builder{}
	sep := t.TabSeparator
	for i, tab := range t.tabs {
		style := t.ThemeMap[t.Strategy.ThemeKey]
		if i == t.activeTab {
			style = render.Theme{
				TextColor:   "#fff",
				BgColor:     "#333",
				BorderColor: "#f90",
			}
		}
		prefix := "  "
		if t.UseDot && i == t.activeTab {
			prefix = "• "
		}
		if t.UseDot {
			s.WriteString(prefix)
		}
		styled := render.ApplyStrategy(tab, t.Strategy, style)
		s.WriteString(styled)
		if i != len(t.tabs)-1 {
			s.WriteString(sep)
		}
	}
	if t.Width > 0 {
		return render.ApplyStrategy(s.String(), t.Strategy, t.ThemeMap[t.Strategy.ThemeKey])
	}
	return s.String()
}

func (t *Tabs) activeTabCmd() tea.Msg {
	return ActiveTabMsg(t.activeTab)
}

// SelectTabCmd is a Bubble Tea command that selects the tab at the given index.
func SelectTabCmd(tab int) tea.Cmd {
	return func() tea.Msg {
		return SelectTabMsg(tab)
	}
}

// Exported getter for the active tab index
func (t *Tabs) ActiveTab() int {
	return t.activeTab
}

// Exported getter for the tab names
func (t *Tabs) TabNames() []string {
	return t.tabs
}
