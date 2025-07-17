package chat
package chat

import (
	"aichat/interfaces"
	"aichat/types"
	render "aichat/types/render"
)

// Factory functions for chat-related ViewState implementations

func NewCompositeChatViewStateFactory(
	ctx types.Context,
	nav interfaces.Controller,
	tabs []string,
	recent []string,
	favorites []string,
	themeMap render.ThemeMap,
	strategies map[string]render.RenderStrategy,
) *CompositeChatViewState {
	return NewCompositeChatViewState(ctx, nav, tabs, recent, favorites, themeMap, strategies)
}

func NewSidebarTopModalFactory() *SidebarTopModal {
	return NewSidebarTopModal()
}

func NewSidebarBottomModalFactory() *SidebarBottomModal {
	return NewSidebarBottomModal()
}

func NewChatWindowModalFactory() *ChatWindowModal {
	return NewChatWindowModal()
}

func NewInputAreaModalFactory(onSubmit func(func(string) error)) *InputAreaModal {
	return NewInputAreaModal(onSubmit)
}

// --- IMPLEMENTATIONS FOR CHAT VIEWSTATE AND MODALS ---

// CompositeChatViewState is the main container for the chat UI layout.
type CompositeChatViewState struct {
	Ctx        types.Context
	Nav        interfaces.Controller
	Tabs       []string
	Recent     []string
	Favorites  []string
	ThemeMap   render.ThemeMap
	Strategies map[string]render.RenderStrategy
}

func NewCompositeChatViewState(ctx types.Context, nav interfaces.Controller, tabs []string, recent []string, favorites []string, themeMap render.ThemeMap, strategies map[string]render.RenderStrategy) *CompositeChatViewState {
	return &CompositeChatViewState{
		Ctx:        ctx,
		Nav:        nav,
		Tabs:       tabs,
		Recent:     recent,
		Favorites:  favorites,
		ThemeMap:   themeMap,
		Strategies: strategies,
	}
}

func (c *CompositeChatViewState) View() string {
	return "[CompositeChatViewState: main chat UI layout]"
}

// SidebarTopModal represents a modal in the top sidebar region.
type SidebarTopModal struct {
	Title   string
	Content string
}

func NewSidebarTopModal() *SidebarTopModal {
	return &SidebarTopModal{
		Title:   "Sidebar Top Modal",
		Content: "[SidebarTopModal content]",
	}
}

func (m *SidebarTopModal) View() string {
	return m.Title + ": " + m.Content
}

// SidebarBottomModal represents a modal in the bottom sidebar region.
type SidebarBottomModal struct {
	Title   string
	Content string
}

func NewSidebarBottomModal() *SidebarBottomModal {
	return &SidebarBottomModal{
		Title:   "Sidebar Bottom Modal",
		Content: "[SidebarBottomModal content]",
	}
}

func (m *SidebarBottomModal) View() string {
	return m.Title + ": " + m.Content
}

// ChatWindowModal represents a modal in the chat window region.
type ChatWindowModal struct {
	Title   string
	Content string
}

func NewChatWindowModal() *ChatWindowModal {
	return &ChatWindowModal{
		Title:   "Chat Window Modal",
		Content: "[ChatWindowModal content]",
	}
}

func (m *ChatWindowModal) View() string {
	return m.Title + ": " + m.Content
}

// InputAreaModal represents a modal in the input area region.
type InputAreaModal struct {
	Title    string
	Content  string
	OnSubmit func(func(string) error)
}

func NewInputAreaModal(onSubmit func(func(string) error)) *InputAreaModal {
	return &InputAreaModal{
		Title:    "Input Area Modal",
		Content:  "[InputAreaModal content]",
		OnSubmit: onSubmit,
	}
}

func (m *InputAreaModal) View() string {
	status := "[no submit handler]"
	if m.OnSubmit != nil {
		status = "[submit handler set]"
	}
	return m.Title + ": " + m.Content + " " + status
}
