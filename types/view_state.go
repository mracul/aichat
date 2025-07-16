// types/view_state.go - Unified ViewState implementations for navigation stack
// MIGRATION TARGET: legacy/gui.go (all menu, chat, modal state logic)

package types

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewState interface: all view states (menus, modals, chats, etc.) implement this
// OOP: Enables polymorphism and stack-based navigation
// Design: Context and controller are injected for testability and decoupling

type ViewState interface {
	Type() ViewType
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd)
	View() string
	Init() tea.Cmd
	// Added for navigation.go compatibility
	IsMainMenu() bool
	MarshalState() ([]byte, error)
	ViewType() ViewType
	UnmarshalState([]byte) error
}

// MenuViewState: represents a menu or submenu view
// Holds a MenuEntrySet and cursor, and is context/controller-driven

type MenuViewState struct {
	menuType     MenuType
	entries      MenuEntrySet
	cursor       int
	title        string
	ctx          Context
	nav          Controller
	WindowWidth  int // new: window width in cells
	WindowHeight int // new: window height in cells
}

func NewMenuViewState(menuType MenuType, entries MenuEntrySet, title string, ctx Context, nav Controller) *MenuViewState {
	return &MenuViewState{
		menuType: menuType,
		entries:  entries,
		cursor:   0,
		title:    title,
		ctx:      ctx,
		nav:      nav,
	}
}

func (mvs *MenuViewState) Type() ViewType { return MenuStateType }

func (mvs *MenuViewState) Init() tea.Cmd { return nil }

func (mvs *MenuViewState) View() string {
	// ASCII art heading
	ascii := ""
	if mvs.WindowWidth > 0 {
		ascii = centerString(menuAsciiArt(), mvs.WindowWidth)
	} else {
		ascii = menuAsciiArt()
	}
	// Title
	titleColor := MenuTitleColorMap[mvs.menuType]
	titleStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).MarginBottom(1)
	if titleColor != "" {
		titleStyle = titleStyle.Foreground(lipgloss.Color(titleColor))
	}
	title := titleStyle.Render(mvs.title)
	if mvs.WindowWidth > 0 {
		title = lipgloss.PlaceHorizontal(mvs.WindowWidth, lipgloss.Center, title)
	}
	// Dynamic menu box sizing
	menuWidth := 300
	menuHeight := 400
	if mvs.WindowWidth > 0 {
		w := int(float64(mvs.WindowWidth) * 0.382)
		if w > menuWidth {
			menuWidth = w
		}
	}
	if mvs.WindowHeight > 0 {
		h := int(float64(mvs.WindowHeight) * 0.618)
		if h > menuHeight {
			menuHeight = h
		}
	}
	// Entries
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
		if entry.Description != "" {
			text += " — " + entry.Description
		}
		rendered = append(rendered, style.Render(text))
	}
	menuBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Width(menuWidth).Height(menuHeight).Align(lipgloss.Center).Render(lipgloss.JoinVertical(lipgloss.Left, rendered...))
	if mvs.WindowWidth > 0 {
		menuBox = lipgloss.PlaceHorizontal(mvs.WindowWidth, lipgloss.Center, menuBox)
	}
	// Control hints (left-aligned, beneath menu box)
	meta := MenuMetas[mvs.menuType]
	controls := ControlInfoMap[meta.ControlInfoType]
	var controlLines []string
	for _, line := range controls.Lines {
		controlLines = append(controlLines, lipgloss.NewStyle().Faint(true).Render(line))
	}
	controlInfo := lipgloss.JoinVertical(lipgloss.Left, controlLines...)
	// Compose: ASCII heading, title, menu box, control info (beneath)
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		ascii,
		title,
		menuBox,
		controlInfo,
	)
	w, h := mvs.WindowWidth, mvs.WindowHeight
	if w == 0 {
		w = 800
	}
	if h == 0 {
		h = 600
	}
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

// Bubble Tea compatibility: wraps OOP-style Update
func (mvs *MenuViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return mvs.UpdateWithContext(msg, mvs.ctx, mvs.nav)
}

// OOP-style update
func (mvs *MenuViewState) UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if mvs.cursor > 0 {
				mvs.cursor--
			}
		case "down", "j":
			if mvs.cursor < len(mvs.entries)-1 {
				mvs.cursor++
			}
		case "enter":
			entry := mvs.entries[mvs.cursor]
			if entry.Disabled {
				return mvs, nil
			}
			// Always call Action if present
			if entry.Action != nil {
				if err := entry.Action(ctx, nav); err != nil {
					nav.ShowModal("error", err.Error())
				}
				// If Action handled navigation, do not also push Next
				return mvs, nil
			}
			// If no Action, but Next is set, navigate to the next menu
			if entry.Next != 0 {
				nav.Push(NewMenuViewState(entry.Next, getMenuEntries(entry.Next), menuTypeToString(entry.Next), ctx, nav))
			}
		case "esc", "q":
			if nav.CanPop() {
				nav.Pop()
			}
		}
	}
	return mvs, nil
}

func (mvs *MenuViewState) IsMainMenu() bool              { return mvs.menuType == MainMenu }
func (mvs *MenuViewState) MarshalState() ([]byte, error) { return nil, nil }
func (mvs *MenuViewState) ViewType() ViewType            { return MenuStateType }
func (mvs *MenuViewState) UnmarshalState([]byte) error   { return nil }

func (mvs *MenuViewState) Resize(w, h int) {
	if w < 800 {
		w = 800
	}
	if h < 600 {
		h = 600
	}
	mvs.WindowWidth = w
	mvs.WindowHeight = h
}

// ChatViewState: compatible with tea.Model and OOP-style update
// (Stub for demonstration)
type ChatViewState struct {
	ChatTitle string
	Messages  []string // TODO: Replace with message structs
	Streaming bool
	ctx       Context
	nav       Controller
}

func (c *ChatViewState) Type() ViewType { return ChatStateType }
func (c *ChatViewState) Init() tea.Cmd  { return nil }
func (c *ChatViewState) View() string   { return "[Chat: " + c.ChatTitle + "]" }
func (c *ChatViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.UpdateWithContext(msg, c.ctx, c.nav)
}
func (c *ChatViewState) UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd) {
	// TODO: Implement chat update logic
	return c, nil
}

func (c *ChatViewState) IsMainMenu() bool              { return false }
func (c *ChatViewState) MarshalState() ([]byte, error) { return nil, nil }
func (c *ChatViewState) ViewType() ViewType            { return ChatStateType }
func (c *ChatViewState) UnmarshalState([]byte) error   { return nil }

// ModalViewState: compatible with tea.Model and OOP-style update
// (Stub for demonstration)
type ModalViewState struct {
	ModalType string
	Content   string
	ctx       Context
	nav       Controller
}

func (m *ModalViewState) Type() ViewType { return ModalStateType }
func (m *ModalViewState) Init() tea.Cmd  { return nil }
func (m *ModalViewState) View() string   { return "[Modal: " + m.ModalType + "] " + m.Content }
func (m *ModalViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.UpdateWithContext(msg, m.ctx, m.nav)
}
func (m *ModalViewState) UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd) {
	// TODO: Implement modal update logic
	return m, nil
}

func (m *ModalViewState) IsMainMenu() bool              { return false }
func (m *ModalViewState) MarshalState() ([]byte, error) { return nil, nil }
func (m *ModalViewState) ViewType() ViewType            { return ModalStateType }
func (m *ModalViewState) UnmarshalState([]byte) error   { return nil }

// Helper: get menu entries for a given MenuType
func getMenuEntries(menuType MenuType) MenuEntrySet {
	switch menuType {
	case MainMenu:
		return MainMenuEntries
	case ChatsMenu:
		return ChatsMenuEntries
	case FavoritesMenu:
		return FavoritesMenuEntries
	case PromptsMenu:
		return PromptsMenuEntries
	case ModelsMenu:
		return ModelsMenuEntries
	case APIKeyMenu:
		return APIKeyMenuEntries
	case HelpMenu:
		return HelpMenuEntries
	case ExitMenu:
		return ExitMenuEntries
	case SettingsMenu:
		return SettingsMenuEntries
	case ProvidersMenu:
		return ProvidersMenuEntries
	case ThemesMenu:
		return ThemesMenuEntries
	default:
		return nil
	}
}

// QuitAppMsg is sent when the user confirms quitting the app
// Used by quit confirmation modal

type QuitAppMsg struct{}

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
