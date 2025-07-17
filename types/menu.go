package types
package types

import (
	"aichat/interfaces"
)

type MenuType int

const (
	MainMenu MenuType = iota
	ChatsMenu
	FavoritesMenu
	PromptsMenu
	ModelsMenu
	APIKeyMenu
	HelpMenu
	ExitMenu
	SettingsMenu
	ProvidersMenu // Added for settings submenu
	ThemesMenu    // Added for settings submenu
)

type MenuAction func(ctx interfaces.Context, nav interfaces.Controller) error

type MenuEntry struct {
	Text        string
	Description string
	Action      MenuAction
	Next        MenuType
	Disabled    bool
	Shortcut    string
}

type MenuEntrySet []MenuEntry

