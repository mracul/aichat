package sidebar

// SidebarView renders the horizontal chat tabs using SidebarTabsModel.
func (m *SidebarTabsModel) SidebarView() string {
	return m.View()
}
