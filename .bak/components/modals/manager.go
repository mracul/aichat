package modals
// manager.go - Contains ModalManager for managing a stack of modals and restoring previous state on close.

package modals

import (
	"aichat/types"
)

// ModalManager manages a stack of modals for the application and handles observer registration/unregistration.
type ModalManager struct {
	stack   []types.ViewState
	Subject types.Subject // The model/state to observe (e.g., chat, flow, etc.)
}

// Push adds a modal to the top of the stack and registers it as an observer if applicable.
func (mm *ModalManager) Push(modal types.ViewState) {
	if mm.Subject != nil {
		if observer, ok := modal.(types.Observer); ok {
			mm.Subject.RegisterObserver(observer)
		}
	}
	mm.stack = append(mm.stack, modal)
}

// Pop removes and returns the top modal from the stack and unregisters it as an observer if applicable.
func (mm *ModalManager) Pop() types.ViewState {
	if len(mm.stack) == 0 {
		return nil
	}
	m := mm.stack[len(mm.stack)-1]
	mm.stack = mm.stack[:len(mm.stack)-1]
	if mm.Subject != nil {
		if observer, ok := m.(types.Observer); ok {
			mm.Subject.UnregisterObserver(observer)
		}
	}
	return m
}

// Current returns the top modal on the stack without removing it.
func (mm *ModalManager) Current() types.ViewState {
	if len(mm.stack) == 0 {
		return nil
	}
	return mm.stack[len(mm.stack)-1]
}
