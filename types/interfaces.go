package types
package types

import (
	"aichat/interfaces"
)

// ViewType is an enum for different view state types.
type ViewType = interfaces.ViewType

const (
	MenuStateType ViewType = iota
	ChatStateType
	ModalStateType
)

// ViewState interface is defined in view_state.go
// This file contains other interfaces and types

// Observer pattern interfaces
// Observer receives notifications of state changes
type Observer interface {
	Notify(event interface{})
}

// Subject manages observers and notifies them
type Subject interface {
	RegisterObserver(o Observer)
	UnregisterObserver(o Observer)
	NotifyObservers(event interface{})
}

// Command pattern interface
// Command encapsulates an action/event
type Command interface {
	Execute(ctx Context, nav interfaces.Controller) error
}

// NavigationController defines navigation stack methods
type NavigationController interface {
	Push(view ViewState)
	Pop() ViewState
	Replace(view ViewState)
	Current() ViewState
	CanPop() bool
}

// ModalController defines modal management
type ModalController interface {
	ShowModal(modalType ModalType, data interface{})
	HideModal()
}

