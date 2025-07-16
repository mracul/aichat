// composite.go - CompositeChatViewState and related logic for orchestrating chat regions.
// Cleans up commented-out code and centralizes Lipgloss styles for input prompt modal.

package chat

import (
	"aichat/components/modals/dialogs"
	"aichat/components/sidebar"
	"aichat/services/storage/repositories"
	"aichat/types"
	"aichat/types/flows"

	"fmt"

	"time"

	"context"

	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Centralized styles for input prompt modal
var (
	inputPromptStyle = lipgloss.NewStyle().Align(lipgloss.Center).MarginBottom(1)
	inputBoxStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 2).Width(32)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Align(lipgloss.Center).MarginTop(1)
	controlHintStyle = lipgloss.NewStyle().Align(lipgloss.Center).MarginTop(1).Foreground(lipgloss.Color("244"))
)

// View renders the input prompt modal.
func (m *InputPromptModal) View() string {
	prompt := inputPromptStyle.Render(m.InstructionText)
	input := m.Value
	if m.Size > 1 {
		input = lipgloss.NewStyle().Width(32).MaxWidth(32).Render(wrapText(m.Value, 32))
	}
	inputBox := inputBoxStyle.Render(input)
	errorMsg := ""
	if m.Error != "" {
		errorMsg = errorStyle.Render(m.Error)
	}
	control := controlHintStyle.Render(m.ControlText)
	return lipgloss.Place(60, 12, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, prompt, inputBox, errorMsg, control),
	)
}

// Update handles text input and control events for the input prompt modal.
func (m *InputPromptModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement text input, Enter, Esc, and multi-line support
	return m, nil
}

// Add a Type() types.ViewType method to InputPromptModal to satisfy types.ViewState
func (m *InputPromptModal) Type() types.ViewType { return types.ModalStateType }

// Add a stub UpdateWithContext method to InputPromptModal to satisfy types.ViewState
func (m *InputPromptModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	return m, nil
}

// Refactor MenuViewState to use MenuEntrySet
// Replace Entries []string with Entries types.MenuEntrySet
// Add logic to handle Action and Next for each entry
// Implement stubs for all menu actions (List, Add, Remove, Set Default, etc.)
// Wire up navigation using Next and modal/dialog system
// Use Bubble Tea and lipgloss for rendering and input handling
type MenuViewState struct {
	MenuName    string
	Entries     types.MenuEntrySet
	Selected    int
	ControlInfo string
}

// IsMainMenu returns true if this is the main menu
func (m *MenuViewState) IsMainMenu() bool { return m.MenuName == "Main Menu" }

// MarshalState serializes the menu state
func (m *MenuViewState) MarshalState() ([]byte, error) { return nil, nil }

// UnmarshalState deserializes the menu state
func (m *MenuViewState) UnmarshalState(data []byte) error { return nil }

// ViewType returns ModalStateType for menus
func (m *MenuViewState) ViewType() types.ViewType { return types.ModalStateType }

// Init initializes the first modal in the flow
func (m *MenuViewState) Init() tea.Cmd { return nil }

// Update method: handle up/down/enter/esc, return new MenuViewState for submenus
func (m *MenuViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up":
			if m.Selected > 0 {
				m.Selected--
			} else {
				m.Selected = len(m.Entries) - 1
			}
			return m, nil
		case "down":
			if m.Selected < len(m.Entries)-1 {
				m.Selected++
			} else {
				m.Selected = 0
			}
			return m, nil
		case "enter":
			// Example: handle Chats submenu
			if m.MenuName == "Main Menu" {
				selected := m.Entries[m.Selected]
				if selected.Text == "Chats" {
					chatsMenu := &MenuViewState{
						MenuName: "Chats",
						Entries: types.MenuEntrySet{
							{Text: "New Chat", Description: "Create a new chat"},
							{Text: "List Chats", Description: "View all chats"},
							{Text: "List Favorites", Description: "View favorite chats"},
							{Text: "Custom Chat", Description: "Create a chat with custom prompts"},
							{Text: "Delete Chat", Description: "Delete a chat"},
						},
						Selected:    0,
						ControlInfo: "Up/Down: Navigate  Enter: Select  Esc: Back",
					}
					return chatsMenu, nil
				}
				// TODO: handle other main menu entries
			}
			if m.MenuName == "Chats" {
				// TODO: handle chats submenu actions
				return m, nil
			}
			return m, nil
		case "esc":
			// Signal back/exit by returning nil (caller should handle stack)
			return nil, nil
		}
	}
	return m, nil
}

// RenderMenuView renders the current menu or submenu using RenderMenuSubmenu
func (m *MenuViewState) RenderMenuView(width, height int) string {
	styledEntries := make([]string, len(m.Entries))
	for i, entry := range m.Entries {
		entryStyle := lipgloss.NewStyle().Padding(0, 1)
		if i == m.Selected {
			entryStyle = entryStyle.Bold(true).Foreground(lipgloss.Color("203")).Background(lipgloss.Color("236"))
		}
		styledEntries[i] = entryStyle.Render(entry.Text)
	}
	// If RenderMenuSubmenu is undefined, comment out or replace with a placeholder
	// menuContent := RenderMenuSubmenu(m.MenuName, styledEntries, m.ControlInfo)
	menuContent := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(m.MenuName),
		lipgloss.JoinVertical(lipgloss.Center, styledEntries...),
		controlHintStyle.Render(m.ControlInfo),
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, menuContent)
}

// CompositeChatViewState is the root controller for the modular chat view
// It implements types.ViewState and manages all four regions

type CompositeChatViewState struct {
	Sidebar      *sidebar.SidebarTabsModel
	ChatWindow   types.ViewState
	InputArea    types.ViewState
	Focused      RegionType
	Layout       LayoutState
	Context      types.Context
	Nav          types.Controller
	Chats        map[string]*ChatViewState
	ActiveChatID string
}

// NewCompositeChatViewState constructs a new composite chat view with the new sidebar
func NewCompositeChatViewState(ctx types.Context, nav types.Controller, tabs, recent, favorites []string) *CompositeChatViewState {
	return &CompositeChatViewState{
		Sidebar:      sidebar.NewSidebarTabsModel(tabs),
		ChatWindow:   NewChatWindowModal(),
		InputArea:    NewInputAreaModal(),
		Focused:      SidebarTop, // or appropriate initial focus
		Layout:       LayoutState{},
		Context:      ctx,
		Nav:          nav,
		Chats:        make(map[string]*ChatViewState),
		ActiveChatID: "",
	}
}

// Type returns the view type
func (cv *CompositeChatViewState) Type() string {
	return "CompositeChatView"
}

// Init initializes all regions
func (cv *CompositeChatViewState) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	// for _, region := range []types.ViewState{cv.Sidebar, cv.ChatWindow, cv.InputArea} {
	for _, region := range []types.ViewState{cv.ChatWindow, cv.InputArea} {
		cmds = append(cmds, region.Init())
	}
	return tea.Batch(cmds...)
}

// View renders the composite chat view, including the new sidebar
func (cv *CompositeChatViewState) View() string {
	// sidebarView := cv.Sidebar.View()
	chatView := cv.ChatWindow.View()
	inputView := cv.InputArea.View()
	// return lipgloss.JoinHorizontal(lipgloss.Left, sidebarView, chatView) + "\n" + inputView
	return chatView + "\n" + inputView
}

// UpdateWithContext delegates updates to the focused region
func (cv *CompositeChatViewState) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// if cv.Focused == SidebarTop {
	// 	model, cmd := cv.Sidebar.Update(msg)
	// 	cv.Sidebar = model.(sidebar.ChatTabsSidebarModel)
	// 	// Detect Enter key
	// 	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
	// 		selected := cv.Sidebar.ActiveChat()
	// 		if selected == "[+] New Chat" {
	// 			// Launch new chat creation in a goroutine, send result as a message
	// 			return cv, func() tea.Msg {
	// 				if name, err := promptForNewChatName(); err == nil && name != "" {
	// 					return newChatCreatedMsg{Name: name}
	// 				}
	// 				return nil
	// 			}
	// 		} else if selected != "" {
	// 			// Switch to selected chat (implement as needed)
	// 			cv.switchToChat(selected)
	// 		}
	// 	}
	// 	// Handle new chat created message
	// 	if newChat, ok := msg.(newChatCreatedMsg); ok && newChat.Name != "" {
	// 		cv.Sidebar = cv.Sidebar.AddTab(newChat.Name)
	// 		cv.switchToChat(newChat.Name)
	// 	}
	// 	return cv, cmd
	// }
	switch cv.Focused {
	// case SidebarTop:
	// 	model, cmd := cv.Sidebar.UpdateWithContext(msg, ctx, nav)
	// 	cv.Sidebar = model.(sidebar.ChatTabsSidebarModel)
	// 	return cv, cmd
	case ChatWindow:
		model, cmd := cv.ChatWindow.UpdateWithContext(msg, ctx, nav)
		cv.ChatWindow = model.(types.ViewState)
		return cv, cmd
	case InputArea:
		model, cmd := cv.InputArea.UpdateWithContext(msg, ctx, nav)
		cv.InputArea = model.(types.ViewState)
		return cv, cmd
	default:
		// Fallback: delegate to chat window
		model, cmd := cv.ChatWindow.UpdateWithContext(msg, ctx, nav)
		cv.ChatWindow = model.(types.ViewState)
		return cv, cmd
	}
}

// Add a stub Update method to CompositeChatViewState to satisfy tea.Model
func (cv *CompositeChatViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return cv, nil }

// Placeholder for control info/help text
func renderControlInfo() string {
	return "Control info help text - F1 to display control list"
}

// --- Region skeletons for composite chat view ---

// SidebarTopModal displays active chats (paginated, selectable)
type SidebarTopModal struct {
	Chats   []string // Placeholder for chat session data
	Cursor  int
	Page    int
	PerPage int
}

func NewSidebarTopModal() *SidebarTopModal {
	return &SidebarTopModal{
		Chats:   []string{"Chat 1", "Chat 2", "Chat 3"},
		Cursor:  0,
		Page:    0,
		PerPage: 10,
	}
}

func (m *SidebarTopModal) Type() types.ViewType                    { return 1002 }
func (m *SidebarTopModal) ViewType() types.ViewType                { return 1002 }
func (m *SidebarTopModal) Init() tea.Cmd                           { return nil }
func (m *SidebarTopModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *SidebarTopModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *SidebarTopModal) MarshalState() ([]byte, error) { return nil, nil }
func (m *SidebarTopModal) UnmarshalState([]byte) error   { return nil }
func (m *SidebarTopModal) IsMainMenu() bool              { return false }
func (m *SidebarTopModal) View() string {
	out := "[Active Chats]"
	for i, chat := range m.Chats {
		marker := "  "
		if i == m.Cursor {
			marker = "> "
		}
		out += "\n" + marker + chat
	}
	return out
}

// SidebarBottomModal displays recent/favorites tabs
// Each tab has its own list and cursor

type SidebarTab int

const (
	TabRecent SidebarTab = iota
	TabFavorites
)

type SidebarBottomModal struct {
	Tab       SidebarTab
	Recent    []string
	Favorites []string
	Cursor    int
	TabCursor [2]int // Cursor per tab
}

func NewSidebarBottomModal() *SidebarBottomModal {
	return &SidebarBottomModal{
		Tab:       TabRecent,
		Recent:    []string{"Recent 1", "Recent 2"},
		Favorites: []string{"Fav 1"},
		Cursor:    0,
		TabCursor: [2]int{0, 0},
	}
}

func (m *SidebarBottomModal) Type() types.ViewType                    { return 1003 }
func (m *SidebarBottomModal) ViewType() types.ViewType                { return 1003 }
func (m *SidebarBottomModal) Init() tea.Cmd                           { return nil }
func (m *SidebarBottomModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *SidebarBottomModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m *SidebarBottomModal) MarshalState() ([]byte, error) { return nil, nil }
func (m *SidebarBottomModal) UnmarshalState([]byte) error   { return nil }
func (m *SidebarBottomModal) IsMainMenu() bool              { return false }
func (m *SidebarBottomModal) View() string {
	tabName := map[SidebarTab]string{TabRecent: "Recent", TabFavorites: "Favorites"}[m.Tab]
	out := "[" + tabName + "]"
	var items []string
	if m.Tab == TabRecent {
		items = m.Recent
	} else {
		items = m.Favorites
	}
	for i, item := range items {
		marker := "  "
		if i == m.TabCursor[m.Tab] {
			marker = "> "
		}
		out += "\n" + marker + item
	}
	return out
}

// ChatWindowModal displays chat messages (skeleton)
type ChatWindowModal struct {
	Messages []string // Placeholder for chat messages
	Scroll   int
}

func NewChatWindowModal() *ChatWindowModal {
	return &ChatWindowModal{
		Messages: []string{"User: Hello", "AI: Hi!"},
		Scroll:   0,
	}
}

func (m *ChatWindowModal) Type() types.ViewType                    { return types.ChatStateType }
func (m *ChatWindowModal) ViewType() types.ViewType                { return types.ChatStateType }
func (m *ChatWindowModal) Init() tea.Cmd                           { return nil }
func (m *ChatWindowModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *ChatWindowModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// Placeholder for scroll, etc.
	return m, nil
}
func (m *ChatWindowModal) MarshalState() ([]byte, error) { return nil, nil }
func (m *ChatWindowModal) UnmarshalState([]byte) error   { return nil }
func (m *ChatWindowModal) IsMainMenu() bool              { return false }
func (m *ChatWindowModal) View() string {
	out := "[Chat Window]"
	for _, msg := range m.Messages {
		out += "\n" + msg
	}
	return out
}

// InputAreaModal is a skeleton for the advanced text editor

type InputAreaModal struct {
	Input  string
	Cursor int
}

func NewInputAreaModal() *InputAreaModal {
	return &InputAreaModal{
		Input:  "",
		Cursor: 0,
	}
}

func (m *InputAreaModal) Type() types.ViewType                    { return 1001 } // Use int constant if types.InputAreaType is undefined
func (m *InputAreaModal) ViewType() types.ViewType                { return 1001 }
func (m *InputAreaModal) Init() tea.Cmd                           { return nil }
func (m *InputAreaModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *InputAreaModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	// Placeholder for input editing
	return m, nil
}
func (m *InputAreaModal) MarshalState() ([]byte, error) { return nil, nil }
func (m *InputAreaModal) UnmarshalState([]byte) error   { return nil }
func (m *InputAreaModal) IsMainMenu() bool              { return false }
func (m *InputAreaModal) View() string {
	return "[Input] > " + m.Input
}

// DynamicNoticeModal cycles through an array of notices at a set interval
// Used for animated 'Testing...' etc.
type DynamicNoticeModal struct {
	Notices       []string
	Current       int
	Interval      int // in seconds
	TickCount     int
	Message       string
	OnComplete    func(success bool)
	Testing       bool
	Success       bool
	Done          bool
	ResultMessage string
	ResultEmoji   string
}

// IsMainMenu returns false for DynamicNoticeModal
func (m *DynamicNoticeModal) IsMainMenu() bool { return false }

// MarshalState serializes the modal state
func (m *DynamicNoticeModal) MarshalState() ([]byte, error) { return nil, nil }

// UnmarshalState deserializes the modal state
func (m *DynamicNoticeModal) UnmarshalState(data []byte) error { return nil }

// ViewType returns ModalStateType
func (m *DynamicNoticeModal) ViewType() types.ViewType { return types.ModalStateType }

// Add a Type() types.ViewType method to DynamicNoticeModal to satisfy types.ViewState
func (m *DynamicNoticeModal) Type() types.ViewType { return types.ModalStateType }

// Init initializes the dynamic notice modal
func (m *DynamicNoticeModal) Init() tea.Cmd { return nil }

// GetControlSets returns the dynamic notice modal's control sets
func (m *DynamicNoticeModal) GetControlSets() []types.ControlSet {
	controls := []types.ControlSet{
		{
			Controls: []types.ControlType{
				{
					Name: "Esc", Key: tea.KeyEsc,
					Action: func() bool {
						// TODO: handle cancel
						return true
					},
				},
			},
		},
	}
	return controls
}

// Update cycles through notices every Interval seconds using tea.Tick
func (m *DynamicNoticeModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Done {
		// Wait for Enter/Esc to confirm if unsuccessful
		if !m.Success {
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "enter", "esc":
					if m.OnComplete != nil {
						m.OnComplete(false)
					}
				}
			}
		}
		return m, nil
	}
	if !m.Testing {
		return m, nil
	}
	// Cycle through notices every Interval seconds
	m.Current = (m.Current + 1) % len(m.Notices)
	m.Message = m.Notices[m.Current]
	return m, nil
}

// View renders the current notice or result
func (m *DynamicNoticeModal) View() string {
	if m.Done {
		if m.Success {
			return lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Render("✅ " + m.ResultMessage)
		} else {
			return lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Render("❌ " + m.ResultMessage + "\n(Enter or Esc to try another key)")
		}
	}
	return lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Render(m.Message)
}

// Add a stub UpdateWithContext method to DynamicNoticeModal to satisfy types.ViewState
func (m *DynamicNoticeModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	return m, nil
}

// InputPromptModal is a reusable modal for text input with instruction and control text
// Size: 1 for single-line, 3 for multi-line (with wrapping)
type InputPromptModal struct {
	Value           string
	InstructionText string // Text above the input box
	ControlText     string // Hint below the input box
	Size            int    // 1 or 3 lines
	Error           string // Optional error message
	OnSubmit        func(value string)
	OnCancel        func()
}

// IsMainMenu returns false for InputPromptModal
func (m *InputPromptModal) IsMainMenu() bool { return false }

// MarshalState serializes the modal state
func (m *InputPromptModal) MarshalState() ([]byte, error) { return nil, nil }

// UnmarshalState deserializes the modal state
func (m *InputPromptModal) UnmarshalState(data []byte) error { return nil }

// ViewType returns ModalStateType
func (m *InputPromptModal) ViewType() types.ViewType { return types.ModalStateType }

// Init initializes the dynamic notice modal
func (m *InputPromptModal) Init() tea.Cmd { return nil }

// GetControlSets returns the input prompt modal's control sets
func (m *InputPromptModal) GetControlSets() []types.ControlSet {
	controls := []types.ControlSet{
		{
			Controls: []types.ControlType{
				{
					Name: "Enter", Key: tea.KeyEnter,
					Action: func() bool {
						// TODO: handle submit
						return true
					},
				},
				{
					Name: "Esc", Key: tea.KeyEsc,
					Action: func() bool {
						// TODO: handle cancel
						return true
					},
				},
			},
		},
	}
	return controls
}

// wrapText is a simple helper for wrapping text to a given width
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}
	var out string
	for i := 0; i < len(text); i += width {
		end := i + width
		if end > len(text) {
			end = len(text)
		}
		out += text[i:end] + "\n"
	}
	return out
}

// NewCompositeChatView creates a new composite chat view with the given chat as active
func NewCompositeChatView(chat *ChatViewState) *CompositeChatViewState {
	chats := make(map[string]*ChatViewState)
	chats[chat.ChatTitle] = chat // Use ChatTitle as key; replace with unique ID if available
	return &CompositeChatViewState{
		Sidebar:      sidebar.NewSidebarTabsModel([]string{"Chats"}),
		ChatWindow:   NewChatWindowModal(),
		InputArea:    NewInputAreaModal(),
		Focused:      SidebarTop,
		Chats:        chats,
		ActiveChatID: chat.ChatTitle,
	}
}

// CreateMainMenu returns a new main menu MenuViewState
func CreateMainMenu() *MenuViewState {
	return &MenuViewState{
		MenuName: "Main Menu",
		Entries: types.MenuEntrySet{
			{Text: "Chats", Description: "View and manage chats"},
			{Text: "Prompts", Description: "Manage prompt templates"},
			{Text: "Models", Description: "Manage AI models"},
			{Text: "API Keys", Description: "Manage API keys"},
			{Text: "Help", Description: "View help"},
			{Text: "Exit", Description: "Exit application"},
		},
		Selected:    0,
		ControlInfo: "Up/Down: Navigate  Enter: Select  Esc: Back",
	}
}

// View returns a placeholder string for now
var defaultMenuView = "[MenuViewState placeholder view]"

func (m *MenuViewState) View() string { return defaultMenuView }

// GetControlSets returns the menu's control sets
func (m *MenuViewState) GetControlSets() []types.ControlSet {
	controls := []types.ControlSet{
		{
			Controls: []types.ControlType{
				{
					Name: "Up", Key: tea.KeyUp,
					Action: func() bool {
						if m.Selected > 0 {
							m.Selected--
						} else {
							m.Selected = len(m.Entries) - 1
						}
						return true
					},
				},
				{
					Name: "Down", Key: tea.KeyDown,
					Action: func() bool {
						if m.Selected < len(m.Entries)-1 {
							m.Selected++
						} else {
							m.Selected = 0
						}
						return true
					},
				},
				{
					Name: "Enter", Key: tea.KeyEnter,
					Action: func() bool {
						// TODO: handle menu selection
						return true
					},
				},
				{
					Name: "Esc", Key: tea.KeyEsc,
					Action: func() bool {
						// TODO: handle back/cancel
						return true
					},
				},
			},
		},
	}
	return controls
}

func (m *MenuViewState) Type() types.ViewType { return types.MenuStateType }

func (cv *CompositeChatViewState) startCustomChatFlow() {
	// Step 1: Prompt for chat name
	inputModal := &InputPromptModal{
		InstructionText: "Enter name for custom chat:",
		ControlText:     "Enter to confirm, Esc to cancel, leave blank for timestamp",
		Size:            1,
		OnSubmit: func(name string) {
			if name == "" {
				name = time.Now().Format("20060102_150405")
			}
			// Step 2: Prompt selection
			prompts, err := loadPrompts() // []flows.Prompt
			if err != nil || len(prompts) == 0 {
				// Show error modal
				notice := &DynamicNoticeModal{
					Notices:  []string{"No prompts available."},
					Current:  0,
					Interval: 1,
				}
				cv.PushModal(notice)
				return
			}
			promptNames := make([]string, len(prompts))
			for i, p := range prompts {
				promptNames[i] = p.Name
			}
			promptModal := &dialogs.ListModal{
				Title:           "Select Prompt",
				Options:         promptNames,
				InstructionText: "Select prompt:",
				ControlText:     "Up/Down: Navigate, Enter: Select, Esc: Cancel",
				OnSelect: func(promptIdx int) {
					// Step 3: Model selection
					modelsList, _, err := loadModelsWithMostRecent()
					if err != nil || len(modelsList) == 0 {
						notice := &DynamicNoticeModal{
							Notices:  []string{"No models available."},
							Current:  0,
							Interval: 1,
						}
						cv.PushModal(notice)
						return
					}
					modelModal := &dialogs.ListModal{
						Title:           "Select Model",
						Options:         modelsList,
						InstructionText: "Select model:",
						ControlText:     "Up/Down: Navigate, Enter: Select, Esc: Cancel",
						OnSelect: func(modelIdx int) {
							// Step 4: API key selection
							keyRepo := repositories.NewAPIKeyRepository()
							apiKeys, err := keyRepo.GetAll()
							if err != nil || len(apiKeys) == 0 {
								notice := &DynamicNoticeModal{
									Notices:  []string{"No API keys available."},
									Current:  0,
									Interval: 1,
								}
								cv.PushModal(notice)
								return
							}
							keyTitles := make([]string, len(apiKeys))
							for i, k := range apiKeys {
								keyTitles[i] = k.Title
							}
							keyModal := &dialogs.ListModal{
								Title:           "Select API Key",
								Options:         keyTitles,
								InstructionText: "Select API key:",
								ControlText:     "Up/Down: Navigate, Enter: Select, Esc: Cancel",
								OnSelect: func(keyIdx int) {
									// Create and persist chat
									chatFile := types.ChatFile{
										Metadata: types.ChatMetadata{
											Title:     name,
											Model:     modelsList[modelIdx],
											CreatedAt: time.Now(),
										},
										Messages: []types.Message{},
									}
									repo := repositories.NewChatRepository()
									_ = repo.Add(chatFile) // TODO: handle error
									// Open chat view (implement as needed)
									cv.testAPIKeyFlow(prompts[promptIdx].Content, modelsList[modelIdx], apiKeys[keyIdx], func() {
										// API key test successful, proceed to chat view
										cv.PopModal() // Close the API key modal
										// Create and push new composite chat view
										newChat := &ChatViewState{ChatTitle: name}
										cv.Chats[name] = newChat
										cv.ActiveChatID = name
										// TODO: Actually push composite to app navigation stack if needed
									}, func() {
										// API key test failed, allow user to retry
										cv.PopModal()            // Close the API key modal
										cv.startCustomChatFlow() // Restart the flow from the beginning
									})
								},
								CloseSelf: func() { cv.PopModal() },
							}
							cv.PushModal(keyModal)
						},
						CloseSelf: func() { cv.PopModal() },
					}
					cv.PushModal(modelModal)
				},
				CloseSelf: func() { cv.PopModal() },
			}
			cv.PushModal(promptModal)
		},
		OnCancel: func() { cv.PopModal() },
	}
	cv.PushModal(inputModal)
}

func (c *CompositeChatViewState) testAPIKeyFlow(prompt string, model string, apiKey types.APIKey, onSuccess func(), onFailure func()) {
	notices := []string{"Testing.", "Testing. .", "Testing. . ."}
	dynModal := &DynamicNoticeModal{
		Notices:       notices,
		Current:       0,
		Interval:      1,
		Testing:       true,
		Message:       notices[0],
		ResultMessage: "",
		ResultEmoji:   "",
	}
	c.PushModal(dynModal)

	// Start API key test in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resultCh := make(chan bool, 1)
		go func() {
			// Prepare request body
			reqBody := struct {
				Model       string          `json:"model"`
				Messages    []types.Message `json:"messages"`
				Stream      bool            `json:"stream"`
				MaxTokens   int             `json:"max_tokens,omitempty"`
				Temperature float64         `json:"temperature,omitempty"`
			}{
				Model:       model,
				Messages:    []types.Message{{Role: "system", Content: prompt}},
				Stream:      true,
				MaxTokens:   16,
				Temperature: 0.7,
			}
			bodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				resultCh <- false
				return
			}
			req, err := http.NewRequestWithContext(ctx, "POST", apiKey.URL, bytes.NewBuffer(bodyBytes))
			if err != nil {
				resultCh <- false
				return
			}
			req.Header.Set("Authorization", "Bearer "+apiKey.Key)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("HTTP-Referer", "https://github.com/go-ai-cli")
			req.Header.Set("X-Title", "Go AI CLI")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				resultCh <- false
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				resultCh <- false
				return
			}
			// Wait for first response chunk
			buf := make([]byte, 256)
			_, err = resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				resultCh <- false
				return
			}
			resultCh <- true
		}()
		var success bool
		select {
		case success = <-resultCh:
		case <-ctx.Done():
			success = false
		}
		// Update modal on main thread
		if success {
			dynModal.Testing = false
			dynModal.Done = true
			dynModal.Success = true
			dynModal.ResultMessage = "API key test successful!"
			dynModal.ResultEmoji = "✅"
			// Proceed after 1s
			time.Sleep(1 * time.Second)
			if onSuccess != nil {
				onSuccess()
			}
		} else {
			dynModal.Testing = false
			dynModal.Done = true
			dynModal.Success = false
			dynModal.ResultMessage = "unsuccessful, please try another key"
			dynModal.ResultEmoji = "❌"
			// Wait for Enter/Esc to confirm, then call onFailure
			dynModal.OnComplete = func(_ bool) {
				if onFailure != nil {
					onFailure()
				}
			}
		}
	}()
}

// Add this helper at the top-level of the file if not already present
func loadModelsWithMostRecent() ([]string, string, error) {
	path := filepath.Join("src", ".config", "models.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	type Model struct {
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
	}
	var models []Model
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, "", err
	}
	var mostRecentModel string
	var mostRecentName string
	for _, model := range models {
		if model.IsDefault {
			mostRecentModel = model.Name
			mostRecentName = model.Name
			break
		}
	}
	if mostRecentModel == "" {
		// Fallback to the first model if no default is found
		if len(models) > 0 {
			mostRecentModel = models[0].Name
			mostRecentName = models[0].Name
		} else {
			return nil, "", fmt.Errorf("no models found in %s", path)
		}
	}
	return []string{mostRecentModel}, mostRecentName, nil
}

func promptForNewChatName() (string, error) {
	prompt := &survey.Input{
		Message: "Enter chat name:",
		Help:    "Press Esc to cancel.",
	}
	var name string
	err := survey.AskOne(prompt, &name)
	if err != nil {
		// if err == survey.ErrInterrupt {
		// 	return "", nil // User cancelled
		// }
		return "", err
	}
	return name, nil
}

func loadPrompts() ([]flows.Prompt, error) {
	path := filepath.Join("src", ".config", "prompts.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var prompts []flows.Prompt
	if err := json.Unmarshal(data, &prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

func (cv *CompositeChatViewState) switchToChat(chatName string) {
	// This function needs to be implemented to switch the active chat
	// It should update the Sidebar, ChatWindow, and InputArea
	// For now, it just logs the action
	fmt.Printf("Switching to chat: %s\n", chatName)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.SelectTab(chatName)
	// Example: Update ChatWindow
	// cv.ChatWindow = NewChatWindowModal() // This would require a new chat view state
	// Example: Update InputArea
	// cv.InputArea = NewInputAreaModal() // This would require a new input area state
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) PushModal(modal types.ViewState) {
	// This function needs to be implemented to push a modal onto the stack
	// For now, it just logs the action
	fmt.Printf("Pushing modal: %s\n", modal.Type())
	// Example: Add to a slice of modals
	// cv.Modals = append(cv.Modals, modal)
}

func (cv *CompositeChatViewState) PopModal() {
	// This function needs to be implemented to pop a modal from the stack
	// For now, it just logs the action
	fmt.Printf("Popping modal\n")
	// Example: Remove from a slice of modals
	// if len(cv.Modals) > 0 {
	// 	cv.Modals = cv.Modals[:len(cv.Modals)-1]
	// }
}

func (cv *CompositeChatViewState) GetActiveChat() string {
	// This function needs to be implemented to get the active chat name
	// For now, it returns a placeholder
	return "No Active Chat"
}

func (cv *CompositeChatViewState) SetActiveChat(chatName string) {
	// This function needs to be implemented to set the active chat name
	// For now, it just logs the action
	fmt.Printf("Setting active chat to: %s\n", chatName)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.SelectTab(chatName)
	// Example: Update ChatWindow
	// cv.ChatWindow = NewChatWindowModal() // This would require a new chat view state
	// Example: Update InputArea
	// cv.InputArea = NewInputAreaModal() // This would require a new input area state
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetChats() map[string]*ChatViewState {
	// This function needs to be implemented to get all chats
	// For now, it returns a placeholder
	return cv.Chats
}

func (cv *CompositeChatViewState) GetChat(chatName string) *ChatViewState {
	// This function needs to be implemented to get a specific chat
	// For now, it returns a placeholder
	return cv.Chats[chatName]
}

func (cv *CompositeChatViewState) AddChat(chat *ChatViewState) {
	// This function needs to be implemented to add a new chat
	// For now, it just logs the action
	fmt.Printf("Adding chat: %s\n", chat.ChatTitle)
	// Example: Add to the map
	cv.Chats[chat.ChatTitle] = chat
}

func (cv *CompositeChatViewState) RemoveChat(chatName string) {
	// This function needs to be implemented to remove a chat
	// For now, it just logs the action
	fmt.Printf("Removing chat: %s\n", chatName)
	// Example: Remove from the map
	delete(cv.Chats, chatName)
}

func (cv *CompositeChatViewState) GetActiveChatID() string {
	// This function needs to be implemented to get the active chat ID
	// For now, it returns a placeholder
	return cv.ActiveChatID
}

func (cv *CompositeChatViewState) SetActiveChatID(chatID string) {
	// This function needs to be implemented to set the active chat ID
	// For now, it just logs the action
	fmt.Printf("Setting active chat ID to: %s\n", chatID)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.SelectTab(chatID) // Assuming chatID is a tab name
	// Example: Update ChatWindow
	// cv.ChatWindow = NewChatWindowModal() // This would require a new chat view state
	// Example: Update InputArea
	// cv.InputArea = NewInputAreaModal() // This would require a new input area state
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetContext() types.Context {
	// This function needs to be implemented to get the context
	// For now, it returns a placeholder
	return cv.Context
}

func (cv *CompositeChatViewState) SetContext(ctx types.Context) {
	// This function needs to be implemented to set the context
	// For now, it just logs the action
	fmt.Printf("Setting context\n")
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, ctx, cv.Nav) // Assuming Nav is a controller
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, ctx, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, ctx, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetNav() types.Controller {
	// This function needs to be implemented to get the navigation controller
	// For now, it returns a placeholder
	return cv.Nav
}

func (cv *CompositeChatViewState) SetNav(nav types.Controller) {
	// This function needs to be implemented to set the navigation controller
	// For now, it just logs the action
	fmt.Printf("Setting navigation controller\n")
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, nav) // Assuming Context is a context
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetLayout() LayoutState {
	// This function needs to be implemented to get the layout state
	// For now, it returns a placeholder
	return cv.Layout
}

func (cv *CompositeChatViewState) SetLayout(layout LayoutState) {
	// This function needs to be implemented to set the layout state
	// For now, it just logs the action
	fmt.Printf("Setting layout\n")
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetFocused() RegionType {
	// This function needs to be implemented to get the focused region
	// For now, it returns a placeholder
	return cv.Focused
}

func (cv *CompositeChatViewState) SetFocused(focused RegionType) {
	// This function needs to be implemented to set the focused region
	// For now, it just logs the action
	fmt.Printf("Setting focused region to: %d\n", focused)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = focused // Or appropriate focus
}

func (cv *CompositeChatViewState) GetTheme() string {
	// This function needs to be implemented to get the current theme
	// For now, it returns a placeholder
	return "default"
}

func (cv *CompositeChatViewState) SetTheme(theme string) {
	// This function needs to be implemented to set the current theme
	// For now, it just logs the action
	fmt.Printf("Setting theme to: %s\n", theme)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetFont() string {
	// This function needs to be implemented to get the current font
	// For now, it returns a placeholder
	return "default"
}

func (cv *CompositeChatViewState) SetFont(font string) {
	// This function needs to be implemented to set the current font
	// For now, it just logs the action
	fmt.Printf("Setting font to: %s\n", font)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetLanguage() string {
	// This function needs to be implemented to get the current language
	// For now, it returns a placeholder
	return "en"
}

func (cv *CompositeChatViewState) SetLanguage(language string) {
	// This function needs to be implemented to set the current language
	// For now, it just logs the action
	fmt.Printf("Setting language to: %s\n", language)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetAPIKeys() []types.APIKey {
	// This function needs to be implemented to get all API keys
	// For now, it returns a placeholder
	return []types.APIKey{}
}

func (cv *CompositeChatViewState) AddAPIKey(apiKey types.APIKey) {
	// This function needs to be implemented to add an API key
	// For now, it just logs the action
	fmt.Printf("Adding API key: %s\n", apiKey.Title)
	// Example: Add to a slice or repository
	// apiKeys := cv.GetAPIKeys()
	// apiKeys = append(apiKeys, apiKey)
	// cv.SetAPIKeys(apiKeys)
}

func (cv *CompositeChatViewState) RemoveAPIKey(apiKeyTitle string) {
	// This function needs to be implemented to remove an API key
	// For now, it just logs the action
	fmt.Printf("Removing API key: %s\n", apiKeyTitle)
	// Example: Remove from a slice or repository
	// apiKeys := cv.GetAPIKeys()
	// for i, k := range apiKeys {
	// 	if k.Title == apiKeyTitle {
	// 		apiKeys = append(apiKeys[:i], apiKeys[i+1:]...)
	// 		break
	// 	}
	// }
	// cv.SetAPIKeys(apiKeys)
}

func (cv *CompositeChatViewState) GetPrompt(promptName string) *flows.Prompt {
	// This function needs to be implemented to get a specific prompt
	// For now, it returns a placeholder
	return nil
}

func (cv *CompositeChatViewState) AddPrompt(prompt flows.Prompt) {
	// This function needs to be implemented to add a new prompt
	// For now, it just logs the action
	fmt.Printf("Adding prompt: %s\n", prompt.Name)
	// Example: Add to a slice or repository
	// prompts := cv.GetPrompts()
	// prompts = append(prompts, prompt)
	// cv.SetPrompts(prompts)
}

func (cv *CompositeChatViewState) RemovePrompt(promptName string) {
	// This function needs to be implemented to remove a prompt
	// For now, it just logs the action
	fmt.Printf("Removing prompt: %s\n", promptName)
	// Example: Remove from a slice or repository
	// prompts := cv.GetPrompts()
	// for i, p := range prompts {
	// 	if p.Name == promptName {
	// 		prompts = append(prompts[:i], prompts[i+1:]...)
	// 		break
	// 	}
	// }
	// cv.SetPrompts(prompts)
}

func (cv *CompositeChatViewState) GetPrompts() []flows.Prompt {
	// This function needs to be implemented to get all prompts
	// For now, it returns a placeholder
	return []flows.Prompt{}
}

func (cv *CompositeChatViewState) SetPrompts(prompts []flows.Prompt) {
	// This function needs to be implemented to set all prompts
	// For now, it just logs the action
	fmt.Printf("Setting prompts\n")
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "prompts.json")
	// data, err := json.Marshal(prompts)
	// if err != nil {
	// 	fmt.Printf("Error marshalling prompts: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing prompts: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModels() []string {
	// This function needs to be implemented to get all model names
	// For now, it returns a placeholder
	return []string{}
}

func (cv *CompositeChatViewState) SetModels(models []string) {
	// This function needs to be implemented to set all model names
	// For now, it just logs the action
	fmt.Printf("Setting models\n")
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModel(modelName string) string {
	// This function needs to be implemented to get a specific model name
	// For now, it returns a placeholder
	return ""
}

func (cv *CompositeChatViewState) SetModel(modelName string) {
	// This function needs to be implemented to set the current model
	// For now, it just logs the action
	fmt.Printf("Setting model to: %s\n", modelName)
	// Example: Update Sidebar
	// cv.Sidebar = cv.Sidebar.UpdateWithContext(nil, cv.Context, cv.Nav) // Assuming Context and Nav are available
	// Example: Update ChatWindow
	// cv.ChatWindow = cv.ChatWindow.UpdateWithContext(nil, cv.Context, cv.Nav)
	// Example: Update InputArea
	// cv.InputArea = cv.InputArea.UpdateWithContext(nil, cv.Context, cv.Nav)
	// cv.Focused = InputArea // Or appropriate focus
}

func (cv *CompositeChatViewState) GetModelDescription(modelName string) string {
	// This function needs to be implemented to get a model description
	// For now, it returns a placeholder
	return ""
}

func (cv *CompositeChatViewState) SetModelDescription(modelName string, description string) {
	// This function needs to be implemented to set a model description
	// For now, it just logs the action
	fmt.Printf("Setting model description for %s: %s\n", modelName, description)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].Description = description
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelPrice(modelName string) float64 {
	// This function needs to be implemented to get a model price
	// For now, it returns a placeholder
	return 0.0
}

func (cv *CompositeChatViewState) SetModelPrice(modelName string, price float64) {
	// This function needs to be implemented to set a model price
	// For now, it just logs the action
	fmt.Printf("Setting model price for %s: %f\n", modelName, price)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].Price = price
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelTokens(modelName string) int {
	// This function needs to be implemented to get a model token limit
	// For now, it returns a placeholder
	return 0
}

func (cv *CompositeChatViewState) SetModelTokens(modelName string, tokens int) {
	// This function needs to be implemented to set a model token limit
	// For now, it just logs the action
	fmt.Printf("Setting model tokens for %s: %d\n", modelName, tokens)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].Tokens = tokens
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelTemperature(modelName string) float64 {
	// This function needs to be implemented to get a model temperature
	// For now, it returns a placeholder
	return 0.0
}

func (cv *CompositeChatViewState) SetModelTemperature(modelName string, temperature float64) {
	// This function needs to be implemented to set a model temperature
	// For now, it just logs the action
	fmt.Printf("Setting model temperature for %s: %f\n", modelName, temperature)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].Temperature = temperature
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelMaxTokens(modelName string) int {
	// This function needs to be implemented to get a model max tokens
	// For now, it returns a placeholder
	return 0
}

func (cv *CompositeChatViewState) SetModelMaxTokens(modelName string, maxTokens int) {
	// This function needs to be implemented to set a model max tokens
	// For now, it just logs the action
	fmt.Printf("Setting model max tokens for %s: %d\n", modelName, maxTokens)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].MaxTokens = maxTokens
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelStream(modelName string) bool {
	// This function needs to be implemented to get a model stream setting
	// For now, it returns a placeholder
	return false
}

func (cv *CompositeChatViewState) SetModelStream(modelName string, stream bool) {
	// This function needs to be implemented to set a model stream setting
	// For now, it just logs the action
	fmt.Printf("Setting model stream for %s: %t\n", modelName, stream)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].Stream = stream
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}

func (cv *CompositeChatViewState) GetModelAPIKey(modelName string) types.APIKey {
	// This function needs to be implemented to get a model API key
	// For now, it returns a placeholder
	return types.APIKey{}
}

func (cv *CompositeChatViewState) SetModelAPIKey(modelName string, apiKey types.APIKey) {
	// This function needs to be implemented to set a model API key
	// For now, it just logs the action
	fmt.Printf("Setting model API key for %s: %s\n", modelName, apiKey.Title)
	// Example: Update a file or repository
	// path := filepath.Join("src", ".config", "models.json")
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Printf("Error reading models: %v\n", err)
	// 	return
	// }
	// var models []Model
	// if err := json.Unmarshal(data, &models); err != nil {
	// 	fmt.Printf("Error unmarshalling models: %v\n", err)
	// 	return
	// }
	// for i, m := range models {
	// 	if m.Name == modelName {
	// 		models[i].APIKey = apiKey
	// 		break
	// 	}
	// }
	// data, err = json.Marshal(models)
	// if err != nil {
	// 	fmt.Printf("Error marshalling models: %v\n", err)
	// 	return
	// }
	// if err := os.WriteFile(path, data, 0644); err != nil {
	// 	fmt.Printf("Error writing models: %v\n", err)
	// 	return
	// }
}
