// model.go - Unified UnifiedAppModel for all TUI/GUI logic, merging previous UnifiedAppModel and GUIAppModel
// This struct centralizes all navigation, modal, UI, and state management for the application.

package app

import (
	"log/slog"
	"time"

	"aichat/components/modals"
	"aichat/components/sidebar"
	"aichat/services/storage"
	"aichat/state"

	// "aichat/src/types" // Remove types import, not needed for AppConfig

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
