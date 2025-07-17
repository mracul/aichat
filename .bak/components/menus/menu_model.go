package menus
package menus

import "aichat/types"

type MenuModel struct {
	entries     types.MenuEntrySet
	menuType    types.MenuType
	selectedIdx int
}

func NewMenuModel(entries types.MenuEntrySet, menuType types.MenuType) *MenuModel {
	return &MenuModel{
		entries:     entries,
		menuType:    menuType,
		selectedIdx: 0,
	}
}

func (m *MenuModel) GetEntries() []types.MenuEntry      { return []types.MenuEntry(m.entries) }
func (m *MenuModel) GetCurrentMenuType() types.MenuType { return m.menuType }
func (m *MenuModel) GetSelectedIndex() int              { return m.selectedIdx }
func (m *MenuModel) SetSelectedIndex(idx int)           { m.selectedIdx = idx }
