package menus
// ChatMenu.go - Contains all logic for the Chats menu and its flows in the menu system.
// This includes actions for listing chats, favoriting/unfavoriting, renaming, previewing, and launching chats.

package menus

import (
	"aichat/types"
	"aichat/types/render"
)

// MenuState represents the current and previous state for menu/modal navigation.
type MenuState struct {
	Current string
	Prev    string
}

// MenuFlow manages a stack of modals and the parent state for a flow.
type MenuFlow struct {
	Modals      []MenuModal
	ParentState string
}

// MenuModal is a generic interface for all modals in a flow.
type MenuModal interface {
	View() string
	Update(input string) (MenuModal, bool) // returns updated modal and whether to advance
	GetPrev() string
}

// InputPromptModal prompts the user for input (e.g., renaming a chat).
type InputPromptModal struct {
	Prompt      string
	Input       string
	ControlInfo string
	ActionInfo  string
	Prev        string
	ThemeMap    render.ThemeMap
	Strategy    render.RenderStrategy
}

// Observer pattern: implement types.Observer
func (m *InputPromptModal) Notify(event interface{}) {
	if ev, ok := event.(types.Event); ok {
		switch ev.Type {
		case "chat_renamed":
			// Optionally update prompt or close modal if relevant
		case "chat_list_updated":
			// Optionally update modal state if relevant
		}
	}
}

// View renders the input prompt modal using ThemeMap and RenderStrategy.
func (m *InputPromptModal) View() string {
	content := m.Prompt + "\n" + m.Input + "\n" + m.ControlInfo + "\n" + m.ActionInfo
	theme := m.ThemeMap[m.Strategy.ThemeKey]
	return render.ApplyStrategy(content, m.Strategy, theme)
}

// Update handles input for the input prompt modal.
func (m *InputPromptModal) Update(input string) (MenuModal, bool) {
	return m, false // TODO: Handle input, return NoticeModal on Enter, or self on edit
}

// GetPrev returns the previous state for the input prompt modal.
func (m *InputPromptModal) GetPrev() string { return m.Prev }

// NoticeModal displays a notice message (e.g., rename result).
type NoticeModal struct {
	Message  string
	Prev     string
	ThemeMap render.ThemeMap
	Strategy render.RenderStrategy
}

// Observer pattern: implement types.Observer
func (m *NoticeModal) Notify(event interface{}) {
	if ev, ok := event.(types.Event); ok {
		switch ev.Type {
		case "chat_deleted":
			// Optionally update message or close modal if relevant
		case "favorite_toggled":
			// Optionally update message if relevant
		}
	}
}

// View renders the notice modal using ThemeMap and RenderStrategy.
func (m *NoticeModal) View() string {
	theme := m.ThemeMap[m.Strategy.ThemeKey]
	return render.ApplyStrategy(m.Message, m.Strategy, theme)
}

// Update handles input for the notice modal.
func (m *NoticeModal) Update(input string) (MenuModal, bool) {
	return m, true // On Enter or Esc, signal to pop modal
}

// GetPrev returns the previous state for the notice modal.
func (m *NoticeModal) GetPrev() string { return m.Prev }

// ListChatsAction displays the list of chats, allows navigation, selection, favoriting, renaming, and previewing.
func ListChatsAction() error {
	return nil // TODO: Implement chat listing, navigation, and key handling (enter, f, r, p)
}

// FavoriteChatAction toggles favorite status for the selected chat.
func FavoriteChatAction(chatName string) error {
	return nil // TODO: Implement favorite/unfavorite logic
}

// RenameChatAction launches the renaming flow, managing state and modals.
func RenameChatAction(state MenuState, chatName string) MenuState {
	return state // TODO: Implement flow logic for renaming (manage modals, transitions, and state)
}

// PreviewChatAction shows the last 3 messages of the selected chat in a modal popup.
func PreviewChatAction(chatName string) error {
	return nil // TODO: Implement preview logic (open EditorModal with last 3 messages)
}

// loadAllChats loads all chat names and metadata.
func loadAllChats() ([]string, error) {
	return nil, nil // TODO: Implement chat loading logic
}

// loadChatMetadata loads chat metadata.
func loadChatMetadata(chatName string) (*types.ChatMetadata, error) {
	return nil, nil // TODO: Implement metadata loading
}
