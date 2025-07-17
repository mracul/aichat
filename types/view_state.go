package types
package types

import (
	"aichat/interfaces"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Renderable defines the rendering capability
type Renderable interface {
	View() string
}

// Updatable defines update logic
type Updatable interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	UpdateWithContext(msg tea.Msg, ctx interfaces.Context, nav interfaces.Controller) (tea.Model, tea.Cmd)
	Init() tea.Cmd
}

// Serializable defines state (un)marshalling
type Serializable interface {
	MarshalState() ([]byte, error)
	UnmarshalState([]byte) error
}

// Navigable defines navigation-related methods
type Navigable interface {
	IsMainMenu() bool
	Type() interfaces.ViewType
	ViewType() interfaces.ViewType
}

// ViewState composes all the above
type ViewState interface {
	Renderable
	Updatable
	Serializable
	Navigable
}

// ControlSetProvider defines context-specific controls for a view
// Implement for views that provide dynamic controls (e.g., input, chat)
type ControlSetProvider interface {
	GetControlSet(ctx interfaces.Context) ControlSet
}

// MenuViewState: represents a menu or submenu view
// Holds a MenuEntrySet and cursor, and is context/controller-driven

type MenuViewState struct {
	menuType     MenuType
	entries      MenuEntrySet
	cursor       int
	title        string
	ctx          interfaces.Context
	Nav          interfaces.Controller
	WindowWidth  int // new: window width in cells
	WindowHeight int // new: window height in cells
}

func NewMenuViewState(menuType MenuType, entries MenuEntrySet, title string, ctx interfaces.Context, nav interfaces.Controller) *MenuViewState {
	return &MenuViewState{
		menuType: menuType,
		entries:  entries,
		cursor:   0,
		title:    title,
		ctx:      ctx,
		Nav:      nav,
	}
}

func (mvs *MenuViewState) Type() interfaces.ViewType { return interfaces.MenuStateType }

// Exported getter for menuType
func (mvs *MenuViewState) MenuType() MenuType { return mvs.menuType }

// Exported getter for cursor
func (mvs *MenuViewState) Cursor() int { return mvs.cursor }

func (mvs *MenuViewState) Init() tea.Cmd { return nil }

func (mvs *MenuViewState) View() string {
	w, h := mvs.WindowWidth, mvs.WindowHeight
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	// ASCII header (centered horizontally)
	ascii := centerString(menuAsciiArt(), w)

	// Title (centered in menu box)
	titleColor := MenuTitleColorMap[mvs.menuType]
	titleStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
	if titleColor != "" {
		titleStyle = titleStyle.Foreground(lipgloss.Color(titleColor))
	}
	title := titleStyle.Render(mvs.title)

	// Entries (left-aligned, spaced)
	var rendered []string
	for i, entry := range mvs.entries {
		style := lipgloss.NewStyle()
		if i == mvs.cursor {
			style = style.Bold(true).Foreground(lipgloss.Color("203")).Background(lipgloss.Color("236"))
		}
		if entry.Disabled {
			style = style.Faint(true)
		}
		text := entry.Text
		desc := ""
		if entry.Description != "" {
			desc = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("245")).Render(" - " + entry.Description)
		}
		rendered = append(rendered, style.Render(text)+desc)
	}

	// Compose menu content: title centered, then blank line, then left-aligned entries
	menuContentBlock := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, rendered...),
	)

	// Calculate menu box dimensions
	menuWidth := int(float64(w) * 3.0 / 8.0)
	if menuWidth < 40 {
		menuWidth = 40
	}
	contentHeight := len(rendered) + 2 // +2 for title and blank line
	minHeight := contentHeight + 2     // +2 for top/bottom padding
	proportionalHeight := int(float64(h) * 4.0 / 6.0)
	menuHeight := proportionalHeight
	if menuHeight < minHeight {
		menuHeight = minHeight
	}
	if menuHeight > minHeight+6 { // Cap extra space to 6 lines beyond content
		menuHeight = minHeight + 6
	}

	// Vertically center the menu content inside the menu box
	menuContentBlock = lipgloss.Place(menuWidth-4, menuHeight-2, lipgloss.Center, lipgloss.Center, menuContentBlock)

	menuBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1).
		Width(menuWidth).
		Height(menuHeight).
		Align(lipgloss.Left)

	menuBox := menuBoxStyle.Render(menuContentBlock)
	menuBox = lipgloss.PlaceHorizontal(w, lipgloss.Center, menuBox)

	// Control info as a single line, left-aligned below menu box, with margin
	meta := MenuMetas[mvs.menuType]
	controls := ControlInfoMap[meta.ControlInfoType]
	controlLine := ""
	if len(controls.Lines) > 0 {
		controlLine = controls.Lines[0]
	}
	controlInfo := lipgloss.NewStyle().Faint(true).Width(menuWidth).Align(lipgloss.Left).MarginTop(1).Render(controlLine)
	leftPad := (w - menuWidth) / 2
	controlInfoBlock := lipgloss.PlaceHorizontal(w, lipgloss.Left, strings.Repeat(" ", leftPad)+controlInfo)

	// Compose the group as a vertical block (ascii, menuBox, controlInfoBlock)
	group := lipgloss.JoinVertical(lipgloss.Top,
		ascii,
		menuBox,
		controlInfoBlock,
	)

	// Center the group vertically and horizontally in the available space
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, group)
}

// Bubble Tea compatibility: wraps OOP-style Update
func (mvs *MenuViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return mvs.UpdateWithContext(msg, mvs.ctx, mvs.Nav)
}

// OOP-style update
func (mvs *MenuViewState) UpdateWithContext(msg tea.Msg, ctx interfaces.Context, nav interfaces.Controller) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			mvs.MoveCursorUp()
		case "down", "j":
			mvs.MoveCursorDown()
		case "enter":
			entry := mvs.entries[mvs.cursor]
			if entry.Disabled {
				return mvs, nil
			}
			if entry.Action != nil {
				if err := entry.Action(ctx, nav); err != nil {
					nav.ShowModal("error", err.Error())
				}
				return mvs, nil
			}
			if entry.Next != 0 {
				nav.Push(NewMenuViewState(entry.Next, getMenuEntries(entry.Next), menuTypeToString(entry.Next), ctx, nav))
			}
		case "esc":
			if mvs.IsMainMenu() {
				nav.Pop()
			} else if nav.CanPop() {
				nav.Pop()
			}
		case "ctrl+q":
			nav.Pop()
		}
	}
	return mvs, nil
}

// Ensure MenuViewState implements ControlSetProvider
var _ ControlSetProvider = (*MenuViewState)(nil)

// GetControlSet returns the context-specific controls for the menu
func (mvs *MenuViewState) GetControlSet(ctx interfaces.Context) ControlSet {
	// For now, always return DefaultMenuControlSet; can be extended for context
	return DefaultMenuControlSet
}

// menuTypeToString returns a human-readable menu name (local helper)
func menuTypeToString(mt MenuType) string {
	switch mt {
	case MainMenu:
		return "Main Menu"
	case ChatsMenu:
		return "Chats"
	case PromptsMenu:
		return "Prompts"
	case ModelsMenu:
		return "Models"
	case APIKeyMenu:
		return "API Keys"
	case HelpMenu:
		return "Help"
	case ExitMenu:
		return "Exit"
	case SettingsMenu:
		return "Settings"
	case ProvidersMenu:
		return "Providers"
	case ThemesMenu:
		return "Themes"
	default:
		return "Menu"
	}
}

// menuAsciiArt returns the ASCII art for the menu header.
func menuAsciiArt() string {
	return `
  █████╗ ██╗ ██████╗██╗  ██╗ █████╗ ████████╗
 ██╔══██╗██║██╔════╝██║  ██║██╔══██╗╚══██╔══╝
 ███████║██║██║     ███████║███████║   ██║
 ██╔══██║██║██║     ██╔══██║██╔══██║   ██║
 ██║  ██║██║╚██████╗██║  ██║██║  ██║   ██║
 ╚═╝  ╚═╝╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   `
}

// centerString centers a string in a field of the given width.
func centerString(s string, width int) string {
	lines := strings.Split(s, "\n")
	var centered []string
	for _, line := range lines {
		centered = append(centered, lipgloss.PlaceHorizontal(width, lipgloss.Center, line))
	}
	return lipgloss.JoinVertical(lipgloss.Center, centered...)
}

func (mvs *MenuViewState) Entries() MenuEntrySet {
	return mvs.entries
}

// Add or update navigation logic to cycle cursor
func (mvs *MenuViewState) MoveCursorUp() {
	if len(mvs.entries) == 0 {
		return
	}
	if mvs.cursor == 0 {
		mvs.cursor = len(mvs.entries) - 1
	} else {
		mvs.cursor--
	}
}

func (mvs *MenuViewState) MoveCursorDown() {
	if len(mvs.entries) == 0 {
		return
	}
	if mvs.cursor == len(mvs.entries)-1 {
		mvs.cursor = 0
	} else {
		mvs.cursor++
	}
}

// Ensure IsMainMenu is exported and present
func (mvs *MenuViewState) IsMainMenu() bool { return mvs.menuType == MainMenu }

// Add MarshalState and UnmarshalState to fully implement ViewState
func (mvs *MenuViewState) MarshalState() ([]byte, error) { return nil, nil }
func (mvs *MenuViewState) UnmarshalState([]byte) error   { return nil }

// Add ViewType to fully implement ViewState
func (mvs *MenuViewState) ViewType() interfaces.ViewType { return interfaces.MenuStateType }

