package navigation
package navigation

import (
	"aichat/types"
	"sync"
)

// NavigationStack is a thread-safe stack of ViewState objects for navigation and handles observer registration/unregistration.
type NavigationStack struct {
	mu      sync.RWMutex
	stack   []types.ViewState
	Subject types.Subject // The model/state to observe (e.g., app state, navigation context)
}

// NewNavigationStack creates a new stack with the main menu as the anchor.
func NewNavigationStack(main types.ViewState) *NavigationStack {
	return &NavigationStack{stack: []types.ViewState{main}}
}

// Push adds a new view to the top of the stack and registers it as an observer if applicable.
func (ns *NavigationStack) Push(v types.ViewState) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if ns.Subject != nil {
		if observer, ok := v.(types.Observer); ok {
			ns.Subject.RegisterObserver(observer)
		}
	}
	ns.stack = append(ns.stack, v)
}

// Pop removes and returns the top view, but never pops the last (main menu). Unregisters observer if applicable.
func (ns *NavigationStack) Pop() types.ViewState {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if len(ns.stack) <= 1 {
		return ns.stack[0]
	}
	v := ns.stack[len(ns.stack)-1]
	ns.stack = ns.stack[:len(ns.stack)-1]
	if ns.Subject != nil {
		if observer, ok := v.(types.Observer); ok {
			ns.Subject.UnregisterObserver(observer)
		}
	}
	return v
}

// ReplaceTop replaces the top view with a new one, handling observer registration/unregistration.
func (ns *NavigationStack) ReplaceTop(v types.ViewState) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	old := ns.stack[len(ns.stack)-1]
	if ns.Subject != nil {
		if observer, ok := old.(types.Observer); ok {
			ns.Subject.UnregisterObserver(observer)
		}
		if observer, ok := v.(types.Observer); ok {
			ns.Subject.RegisterObserver(observer)
		}
	}
	ns.stack[len(ns.stack)-1] = v
}

// Top returns the current top view.
func (ns *NavigationStack) Top() types.ViewState {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.stack[len(ns.stack)-1]
}

// Len returns the number of views on the stack.
func (ns *NavigationStack) Len() int {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return len(ns.stack)
}

// At returns the view at the given index.
func (ns *NavigationStack) At(idx int) types.ViewState {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	if idx < 0 || idx >= len(ns.stack) {
		return nil
	}
	return ns.stack[idx]
}

// MoveToTop moves the view at idx to the top of the stack.
func (ns *NavigationStack) MoveToTop(idx int) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if idx < 0 || idx >= len(ns.stack) {
		return
	}
	v := ns.stack[idx]
	ns.stack = append(ns.stack[:idx], ns.stack[idx+1:]...)
	ns.stack = append(ns.stack, v)
}

// ResetToMainMenu replaces the stack with just the main menu view.
func (ns *NavigationStack) ResetToMainMenu(mainMenu types.ViewState) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.stack = []types.ViewState{mainMenu}
}

