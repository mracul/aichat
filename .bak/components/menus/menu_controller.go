package menus
package menus

import (
	"aichat/components/modals"
	"aichat/components/modals/dialogs"
	"aichat/interfaces"
	"aichat/services/storage"
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
	repo := storage.NewJSONChatRepository("")
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

func NewChatAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("NewChatAction called")
	// TODO: Implement the flow for creating a new chat
	return nil
}

func CustomChatAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("CustomChatAction called")
	// TODO: Implement the flow for creating a custom chat
	return nil
}

func ListPromptsAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("ListPromptsAction called")
	// TODO: Implement the logic to list prompts
	return nil
}

func AddPromptAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("AddPromptAction called")
	// TODO: Implement the logic to add a new prompt
	return nil
}

func SetDefaultPromptAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("SetDefaultPromptAction called")
	// TODO: Implement the logic to set a default prompt
	return nil
}

func DeletePromptAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("DeletePromptAction called")
	// TODO: Implement the logic to delete a prompt
	return nil
}

func ListModelsAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("ListModelsAction called")
	// TODO: Implement the logic to list models
	return nil
}

func AddModelAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("AddModelAction called")
	// TODO: Implement the logic to add a new model
	return nil
}

func SetDefaultModelAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("SetDefaultModelAction called")
	// TODO: Implement the logic to set a default model
	return nil
}

func DeleteModelAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("DeleteModelAction called")
	// TODO: Implement the logic to delete a model
	return nil
}

func ListAPIKeysAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("ListAPIKeysAction called")
	// TODO: Implement the logic to list API keys
	return nil
}

func AddAPIKeyAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("AddAPIKeyAction called")
	// TODO: Implement the logic to add a new API key
	return nil
}

func SetActiveAPIKeyAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("SetActiveAPIKeyAction called")
	// TODO: Implement the logic to set an active API key
	return nil
}

func DeleteAPIKeyAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("DeleteAPIKeyAction called")
	// TODO: Implement the logic to delete an API key
	return nil
}

func ListProvidersAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("ListProvidersAction called")
	// TODO: Implement the logic to list providers
	return nil
}

func AddProviderAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("AddProviderAction called")
	// TODO: Implement the logic to add a new provider
	return nil
}

func ListThemesAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("ListThemesAction called")
	// TODO: Implement the logic to list themes
	return nil
}

func GenerateThemeAction(ctx interfaces.Context, nav interfaces.Controller) error {
	log.Println("GenerateThemeAction called")
	// TODO: Implement the logic to generate a new theme
	return nil
}
