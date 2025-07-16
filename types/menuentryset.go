// menuentryset.go - Defines all menu entry sets for the app's menus and submenus.
// This file centralizes all menu data for extensibility and maintainability.

package types

// Remove incorrect import and use local types

var MainMenuEntries MenuEntrySet
var SettingsMenuEntries MenuEntrySet
var ProvidersMenuEntries MenuEntrySet
var ThemesMenuEntries MenuEntrySet

func init() {
	MainMenuEntries = MenuEntrySet{
		{
			Text:        "Chats",
			Description: "View and manage chats",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(ChatsMenu, getMenuEntries(ChatsMenu), menuTypeToString(ChatsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Prompts",
			Description: "Manage prompt templates",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(PromptsMenu, getMenuEntries(PromptsMenu), menuTypeToString(PromptsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Models",
			Description: "Configure AI models",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(ModelsMenu, getMenuEntries(ModelsMenu), menuTypeToString(ModelsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Help",
			Description: "Show help and shortcuts",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(HelpMenu, getMenuEntries(HelpMenu), menuTypeToString(HelpMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Settings",
			Description: "Configure application settings",
			Action: func(ctx Context, nav Controller) error {
				// Placeholder: implement SettingsMenu and its entries/flow
				nav.Push(NewMenuViewState(SettingsMenu, getMenuEntries(SettingsMenu), menuTypeToString(SettingsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Exit",
			Description: "Exit the application",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(ExitMenu, getMenuEntries(ExitMenu), menuTypeToString(ExitMenu), ctx, nav))
				return nil
			},
		},
	}

	ChatsMenuEntries = MenuEntrySet{
		{
			Text:        "List Chats",
			Description: "View or continue existing chats",
			Action:      nil,
		},
		{
			Text:        "Add New Chat",
			Description: "Create a new chat",
			Action:      nil,
		},
		{
			Text:        "Custom Chat",
			Description: "Custom chat creation (name, model, prompt)",
			Action:      nil,
		},
		{
			Text:        "Load Chat",
			Description: "Load a saved chat",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	FavoritesMenuEntries = MenuEntrySet{
		{
			Text:        "List Favorites",
			Description: "View favorite chats",
			Action:      nil,
		},
		{
			Text:        "Add Favorite",
			Description: "Mark a chat as favorite",
			Action:      nil,
		},
		{
			Text:        "Remove Favorite",
			Description: "Unmark a chat as favorite",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	PromptsMenuEntries = MenuEntrySet{
		{
			Text:        "List Prompts",
			Description: "View all prompts",
			Action:      nil,
		},
		{
			Text:        "Add Prompt",
			Description: "Create a new prompt",
			Action:      nil,
		},
		{
			Text:        "Remove Prompt",
			Description: "Delete a prompt",
			Action:      nil,
		},
		{
			Text:        "Set Default Prompt",
			Description: "Choose default prompt",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	ModelsMenuEntries = MenuEntrySet{
		{
			Text:        "List Models",
			Description: "View all models",
			Action:      nil,
		},
		{
			Text:        "Add Model",
			Description: "Add a new model",
			Action:      nil,
		},
		{
			Text:        "Remove Model",
			Description: "Delete a model",
			Action:      nil,
		},
		{
			Text:        "Set Default Model",
			Description: "Choose default model",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	APIKeyMenuEntries = MenuEntrySet{
		{
			Text:        "List API Keys",
			Description: "View all API keys",
			Action:      nil,
		},
		{
			Text:        "Add API Key",
			Description: "Add a new API key",
			Action:      nil,
		},
		{
			Text:        "Remove API Key",
			Description: "Delete an API key",
			Action:      nil,
		},
		{
			Text:        "Set Active API Key",
			Description: "Choose which key is active",
			Action:      nil,
		},
		{
			Text:        "Test Active Key",
			Description: "Test the current key with a model",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	HelpMenuEntries = MenuEntrySet{
		{
			Text:        "Show Controls",
			Description: "Display controls cheat sheet",
			Action:      nil,
		},
		{
			Text:        "Show About",
			Description: "Display about information",
			Action:      nil,
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	ExitMenuEntries = MenuEntrySet{
		{
			Text:        "Confirm Exit",
			Description: "Exit the application",
			Action:      nil,
		},
		{
			Text:        "Cancel",
			Description: "Return to main menu",
			Next:        MainMenu,
			Action:      nil,
		},
	}

	SettingsMenuEntries = MenuEntrySet{
		{
			Text:        "API Keys",
			Description: "Manage API keys",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(APIKeyMenu, getMenuEntries(APIKeyMenu), menuTypeToString(APIKeyMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Providers",
			Description: "Configure AI providers",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(ProvidersMenu, getMenuEntries(ProvidersMenu), menuTypeToString(ProvidersMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Themes",
			Description: "Select application theme",
			Action: func(ctx Context, nav Controller) error {
				nav.Push(NewMenuViewState(ThemesMenu, getMenuEntries(ThemesMenu), menuTypeToString(ThemesMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Back",
			Description: "Return to main menu",
			Action: func(ctx Context, nav Controller) error {
				nav.Pop()
				return nil
			},
		},
	}

	ProvidersMenuEntries = MenuEntrySet{
		{
			Text:        "Back",
			Description: "Return to settings",
			Action: func(ctx Context, nav Controller) error {
				nav.Pop()
				return nil
			},
		},
	}

	ThemesMenuEntries = MenuEntrySet{
		{
			Text:        "Back",
			Description: "Return to settings",
			Action: func(ctx Context, nav Controller) error {
				nav.Pop()
				return nil
			},
		},
	}
}

// ChatsMenuEntries defines the chats submenu options.
var ChatsMenuEntries = MenuEntrySet{
	{Text: "List Chats", Description: "View or continue existing chats"},
	{Text: "Add New Chat", Description: "Create a new chat"},
	{Text: "Custom Chat", Description: "Custom chat creation (name, model, prompt)"},
	{Text: "Load Chat", Description: "Load a saved chat"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// FavoritesMenuEntries defines the favorites submenu options.
var FavoritesMenuEntries = MenuEntrySet{
	{Text: "List Favorites", Description: "View favorite chats"},
	{Text: "Add Favorite", Description: "Mark a chat as favorite"},
	{Text: "Remove Favorite", Description: "Unmark a chat as favorite"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// PromptsMenuEntries defines the prompts submenu options.
var PromptsMenuEntries = MenuEntrySet{
	{Text: "List Prompts", Description: "View all prompts"},
	{Text: "Add Prompt", Description: "Create a new prompt"},
	{Text: "Remove Prompt", Description: "Delete a prompt"},
	{Text: "Set Default Prompt", Description: "Choose default prompt"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// ModelsMenuEntries defines the models submenu options.
var ModelsMenuEntries = MenuEntrySet{
	{Text: "List Models", Description: "View all models"},
	{Text: "Add Model", Description: "Add a new model"},
	{Text: "Remove Model", Description: "Delete a model"},
	{Text: "Set Default Model", Description: "Choose default model"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// APIKeyMenuEntries defines the API key submenu options.
var APIKeyMenuEntries = MenuEntrySet{
	{Text: "List API Keys", Description: "View all API keys"},
	{Text: "Add API Key", Description: "Add a new API key"},
	{Text: "Remove API Key", Description: "Delete an API key"},
	{Text: "Set Active API Key", Description: "Choose which key is active"},
	{Text: "Test Active Key", Description: "Test the current key with a model"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// HelpMenuEntries defines the help submenu options.
var HelpMenuEntries = MenuEntrySet{
	{Text: "Show Controls", Description: "Display controls cheat sheet"},
	{Text: "Show About", Description: "Display about information"},
	{Text: "Back", Description: "Return to main menu", Next: MainMenu},
}

// ExitMenuEntries defines the exit menu (confirmation modal).
var ExitMenuEntries = MenuEntrySet{
	{Text: "Confirm Exit", Description: "Exit the application"},
	{Text: "Cancel", Description: "Return to main menu", Next: MainMenu},
}
