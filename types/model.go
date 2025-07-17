package types

// model.go - Unified UnifiedAppModel for all TUI/GUI logic, merging previous UnifiedAppModel and GUIAppModel
// This struct centralizes all navigation, modal, UI, and state management for the application.

import (
	"log/slog"
	"time"

	"aichat/components/modals"
	"aichat/components/sidebar"
	"aichat/services/storage"
	"aichat/state"

	// "aichat/src/types" // Remove types import, not needed for AppConfig

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UnifiedAppModel struct {
	// Configuration and logging
	config *AppConfig // Use unqualified type since this is the app package
	logger *slog.Logger

	// Navigation and state management
	navStack *state.NavigationStack
	storage  storage.NavigationStorage

	// UI Components
	sidebar  *sidebar.SidebarTabsModel
	chatView interface{} // Use interface for flexibility

	// Modal system
	modalManager *modals.ModalManager
	modalActive  bool

	// Layout and styling
	width  int
	height int
	style  lipgloss.Style
	styles *GUIStyles // Unified styling for GUI/TUI

	// Performance tracking
	renderCount int64
	lastRender  time.Time

	// Focus and UI state
	focus      string // "main", "modal", "sidebar", "chat", etc.
	isRunning  bool
	helpShown  bool
	showStats  bool
	shouldQuit bool

	// Exit/cleanup
	pendingExitCleanups []func()
}

// NewUnifiedAppModel creates and initializes a new UnifiedAppModel
func NewUnifiedAppModel(config *AppConfig, storage storage.NavigationStorage, logger *slog.Logger) *UnifiedAppModel {
	// Create the main menu view state
	mainMenu := types.NewMenuViewState(
		types.MainMenu,
		types.GetMenuEntries(types.MainMenu),
		"Main Menu",
		nil, // ctx
		nil, // nav (will be set after stack is created)
	)
	// Create the navigation stack with main menu as root
	stack := state.NewNavigationStack(mainMenu)
	// Set nav on mainMenu (circular, but safe for now)
	mainMenuNav := stack
	mainMenuNavIface, _ := any(mainMenu).(*types.MenuViewState)
	if mainMenuNavIface != nil {
		mainMenuNavIface.Nav = mainMenuNav
	}
	return &UnifiedAppModel{
		config:      config,
		logger:      logger,
		storage:     storage,
		navStack:    stack,
		modalActive: false,
		width:       config.DefaultWidth,
		height:      config.DefaultHeight,
		style:       lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")),
		renderCount: 0,
		lastRender:  time.Now(),
		focus:       "navigation",
		isRunning:   true,
		helpShown:   false,
		showStats:   false,
		shouldQuit:  false,
	}
}

// Satisfy tea.Model interface
func (m *UnifiedAppModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *UnifiedAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.navStack != nil && m.navStack.Top() != nil {
			if resizer, ok := m.navStack.Top().(interface{ Resize(int, int) }); ok {
				resizer.Resize(msg.Width, msg.Height)
			}
		}
		return m, nil
	}

	// Forward all other messages to the top view if it implements tea.Model
	if m.navStack != nil && m.navStack.Top() != nil {
		top := m.navStack.Top()
		if teaModel, ok := top.(tea.Model); ok {
			updated, cmd := teaModel.Update(msg)
			// If the returned model is different, update the stack
			if updated != top {
				m.navStack.ReplaceTop(updated.(types.ViewState))
			}
			return m, cmd
		}
	}

	if m.shouldQuit {
		return m, tea.Quit
	}

	return m, nil
}

func (m *UnifiedAppModel) View() string {
	if m.navStack != nil && m.navStack.Top() != nil {
		if vs, ok := m.navStack.Top().(types.ViewState); ok {
			return vs.View()
		}
	}
	return "[No view to render]"
}

// GUIStyles contains all styling for the GUI/TUI
// (Copy from src/gui.go or refactor to a shared location)
type GUIStyles struct {
	mainContainer lipgloss.Style
	headerStyle   lipgloss.Style
	footerStyle   lipgloss.Style
	sidebarStyle  lipgloss.Style
	contentStyle  lipgloss.Style

	titleStyle    lipgloss.Style
	subtitleStyle lipgloss.Style
	textStyle     lipgloss.Style
	helpStyle     lipgloss.Style

	selectedStyle lipgloss.Style
	focusStyle    lipgloss.Style
	disabledStyle lipgloss.Style

	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	warningStyle lipgloss.Style
	infoStyle    lipgloss.Style
}
