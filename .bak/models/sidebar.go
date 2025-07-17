package models
package models

import (
	"aichat/types"
)

type SidebarTabsModel struct {
	Tabs      interface{}
	observers []types.Observer // Observer pattern
}

// Observer pattern methods
func (m *SidebarTabsModel) RegisterObserver(o types.Observer) {
	m.observers = append(m.observers, o)
}

func (m *SidebarTabsModel) UnregisterObserver(o types.Observer) {
	for i, obs := range m.observers {
		if obs == o {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			break
		}
	}
}

func (m *SidebarTabsModel) NotifyObservers(event interface{}) {
	for _, o := range m.observers {
		o.Notify(event)
	}
}

// NewSidebarTabsModel creates a new sidebar tabs model with the given chat names
func NewSidebarTabsModel(tabNames []string) *SidebarTabsModel {
	// Implementation should be in the controller/view layer
	return &SidebarTabsModel{}
}
