package menus

import (
	"aichat/components/modals"
	"aichat/components/modals/dialogs"
	"aichat/interfaces"
	"aichat/services/storage/repositories"
	"aichat/types"
	"log"
)

// MenuController manages menu interactions
type MenuController struct {
	model IMenuModel
	view  IMenuView
}

// NewMenuController creates a new MenuController instance
func NewMenuController(model IMenuModel, view IMenuView) *MenuController {
	return &MenuController{model: model, view: view}
}

// SetView sets the view for the controller
func (c *MenuController) SetView(view IMenuView) {
	c.view = view
}

// SetModel sets the model for the controller
func (c *MenuController) SetModel(model IMenuModel) {
	c.model = model
}

// MoveSelectionUp moves the selection up in the menu
func (c *MenuController) MoveSelectionUp() {
	idx := c.model.GetSelectedIndex()
	if idx > 0 {
		c.model.SetSelectedIndex(idx - 1)
	}
}

// MoveSelectionDown moves the selection down in the menu
func (c *MenuController) MoveSelectionDown() {
	idx := c.model.GetSelectedIndex()
	if idx < len(c.model.GetEntries())-1 {
		c.model.SetSelectedIndex(idx + 1)
	}
}

// Select triggers the action of the selected menu entry
func (c *MenuController) Select(ctx types.Context, nav types.Controller) error {
	entry := c.model.GetEntries()[c.model.GetSelectedIndex()]
	if entry.Action != nil {
		return entry.Action(ctx, nav)
	}
	return nil
}

// ListChatsAction lists chats in a modal with truncated titles
func ListChatsAction(ctx interfaces.Context, nav interfaces.Controller) error {
	repo := repositories.NewChatRepository()
	chats, err := repo.GetAll()
	if err != nil {
		return err
	}
	var chatTitles []string
	for _, chat := range chats {
		title := chat.Metadata.Title
		if len(title) > 20 {
			title = title[:17] + "..."
		}
		chatTitles = append(chatTitles, title)
	}
	// Define what happens when a chat is selected
	onSelect := func(index int) {
		selectedTitle := chatTitles[index]
		// TODO: Implement chat opening logic here
		log.Printf("Selected chat: %s", selectedTitle)
	}
	// Create and push the modal
	modal := dialogs.NewListModalFactory(
		"Chats",
		chatTitles,
		onSelect,
		func() { nav.Pop() },       // closeSelf
		modals.ModalRenderConfig{}, // Use default or pass config
	)
	if nav != nil {
		nav.Push(modal)
	}
	return nil
}

// ... (the rest of the action functions remain the same)
