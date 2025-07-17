package types
// Menu Structure (Canonical)
//
// Main Menu
// ├── Chats
// │   ├── Add new chat (input modal)
// │   ├── List Chats (list view: d=delete, f=favorite, r=rename)
// │   └── Create custom chat (multi-step: name → select prompt → select model)
// ├── Prompts
// │   ├── Add new prompt (input modal - multi step: prompt name then prompt for the text)
// │   ├── Set default prompt (list view)
// │   └── Delete prompt (list view)
// ├── Models
// │   ├── Add model (input modal - multi step: prompt name then prompt for model string)
// │   └── List models (list view: a=set active, d=delete, r=rename)
// ├── Help
// │   ├── Show control overview (modal)
// │   └── Show about (modal)
// ├── Settings
// │   ├── API Keys
// │   │   ├── Add key (input modal multi step - input name, then key, then select provider from list of providers) key stored in schema [name, key, provider, active] json
// │   │   └── Set active key (list view)
// │   ├── Providers
// │   │   └── Add provider (input modal multi step - name then endpoint)
// │   └── Themes
// │       ├── List themes (list view: preview on highlight, set on enter, r rename, d delete)
// │       └── Generate theme (input prompt for name then action)
// └── Exit (confirmation modal)

package types

import (
	"aichat/components/menus"
	"aichat/interfaces"
)

var MainMenuEntries MenuEntrySet
var SettingsMenuEntries MenuEntrySet
var ProvidersMenuEntries MenuEntrySet
var ThemesMenuEntries MenuEntrySet
var ChatsMenuEntries MenuEntrySet
var FavoritesMenuEntries MenuEntrySet
var PromptsMenuEntries MenuEntrySet
var ModelsMenuEntries MenuEntrySet
var APIKeyMenuEntries MenuEntrySet
var HelpMenuEntries MenuEntrySet
var ExitMenuEntries MenuEntrySet

func init() {
	MainMenuEntries = MenuEntrySet{
		{
			Text:        "Chats",
			Description: "View and manage chats",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(ChatsMenu, getMenuEntries(ChatsMenu), menuTypeToString(ChatsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Prompts",
			Description: "Manage prompt templates",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(PromptsMenu, getMenuEntries(PromptsMenu), menuTypeToString(PromptsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Models",
			Description: "Configure AI models",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(ModelsMenu, getMenuEntries(ModelsMenu), menuTypeToString(ModelsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Help",
			Description: "Show help and shortcuts",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(HelpMenu, getMenuEntries(HelpMenu), menuTypeToString(HelpMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Settings",
			Description: "Configure application settings",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(SettingsMenu, getMenuEntries(SettingsMenu), menuTypeToString(SettingsMenu), ctx, nav))
				return nil
			},
		},
		{
			Text:        "Exit",
			Description: "Exit the application",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				// Indirect through high-level handler to avoid import cycle
				if handler, ok := nav.(interface{ HandleExit() }); ok {
					handler.HandleExit()
				}
				return nil
			},
		},
	}

	ChatsMenuEntries = MenuEntrySet{
		{
			Text:        "List Chats",
			Description: "View or continue existing chats",
			Action:      menus.ListChatsAction,
		},
		{
			Text:        "Add New Chat",
			Description: "Create a new chat",
			Action:      menus.NewChatAction,
		},
		{
			Text:        "Create Custom Chat",
			Description: "Multi-step: name → select prompt → select model",
			Action:      menus.CustomChatAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	PromptsMenuEntries = MenuEntrySet{
		{
			Text:   "List Prompts",
			Action: menus.ListPromptsAction,
		},
		{
			Text:   "Add New Prompt",
			Action: menus.AddPromptAction,
		},
		{
			Text:   "Set Default Prompt",
			Action: menus.SetDefaultPromptAction,
		},
		{
			Text:   "Delete Prompt",
			Action: menus.DeletePromptAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	ModelsMenuEntries = MenuEntrySet{
		{
			Text:   "List Models",
			Action: menus.ListModelsAction,
		},
		{
			Text:   "Add Model",
			Action: menus.AddModelAction,
		},
		{
			Text:   "Set Default Model",
			Action: menus.SetDefaultModelAction,
		},
		{
			Text:   "Delete Model",
			Action: menus.DeleteModelAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	SettingsMenuEntries = MenuEntrySet{
		{
			Text: "API Keys",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(APIKeyMenu, getMenuEntries(APIKeyMenu), "API Keys", ctx, nav))
				return nil
			},
		},
		{
			Text: "Providers",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(ProvidersMenu, getMenuEntries(ProvidersMenu), "Providers", ctx, nav))
				return nil
			},
		},
		{
			Text: "Themes",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
				nav.Push(NewMenuViewState(ThemesMenu, getMenuEntries(ThemesMenu), "Themes", ctx, nav))
				return nil
			},
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	APIKeyMenuEntries = MenuEntrySet{
		{
			Text:   "List API Keys",
			Action: menus.ListAPIKeysAction,
		},
		{
			Text:   "Add API Key",
			Action: menus.AddAPIKeyAction,
		},
		{
			Text:   "Set Active API Key",
			Action: menus.SetActiveAPIKeyAction,
		},
		{
			Text:   "Delete API Key",
			Action: menus.DeleteAPIKeyAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	ProvidersMenuEntries = MenuEntrySet{
		{
			Text:   "List Providers",
			Action: menus.ListProvidersAction,
		},
		{
			Text:   "Add Provider",
			Action: menus.AddProviderAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}

	ThemesMenuEntries = MenuEntrySet{
		{
			Text:   "List Themes",
			Action: menus.ListThemesAction,
		},
		{
			Text:   "Generate Theme",
			Action: menus.GenerateThemeAction,
		},
		{
			Text:   "Back",
			Action: func(ctx interfaces.Context, nav interfaces.Controller) error { nav.Pop(); return nil },
		},
	}
}

// Exported function for external use
func GetMenuEntries(menuType MenuType) MenuEntrySet {
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

