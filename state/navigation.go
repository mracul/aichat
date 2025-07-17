package state
// state/navigation.go - NavigationStack using new ViewState protocol
// MIGRATION TARGET: legacy/gui.go (navigation stack logic)

package state

import (
	"aichat/types"
	"encoding/json"
	"sync"
)

// NavigationStack implements LIFO stack with main menu protection
// MIGRATION TARGET: internal/navigation/stack.go

type NavigationStack struct {
	stack []interface{}
	mu    sync.RWMutex
}

// NewNavigationStack creates stack with main menu as root
func NewNavigationStack(mainMenu interface{}) *NavigationStack {
	if mm, ok := mainMenu.(types.ViewState); !ok || !mm.IsMainMenu() {
		panic("root must be main menu")
	}
	return &NavigationStack{stack: []interface{}{mainMenu}}
}

// Push adds state with special handling for main menu
func (ns *NavigationStack) Push(state interface{}) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if mm, ok := state.(types.ViewState); ok && mm.IsMainMenu() {
		ns.stack = []interface{}{state}
		return
	}
	ns.stack = append(ns.stack, state)
}

// Pop removes top state while protecting root
func (ns *NavigationStack) Pop() interface{} {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if len(ns.stack) <= 1 {
		return nil
	}
	top := ns.stack[len(ns.stack)-1]
	ns.stack = ns.stack[:len(ns.stack)-1]
	return top
}

// Top returns current state without popping
func (ns *NavigationStack) Top() interface{} {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	if len(ns.stack) == 0 {
		return nil
	}
	return ns.stack[len(ns.stack)-1]
}

// ReplaceTop replaces the top ViewState with a new one
func (ns *NavigationStack) ReplaceTop(state interface{}) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if len(ns.stack) > 0 {
		ns.stack[len(ns.stack)-1] = state
	}
}

// Replace replaces the top ViewState with a new one
func (ns *NavigationStack) Replace(view interface{}) {
	ns.ReplaceTop(view)
}

// SerializeStack converts entire stack to JSON
func (ns *NavigationStack) SerializeStack() ([]byte, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	var serialized []struct {
		Type  types.ViewType `json:"type"`
		State []byte         `json:"state"`
	}
	for _, s := range ns.stack {
		if vs, ok := s.(types.ViewState); ok {
			stateData, err := vs.MarshalState()
			if err != nil {
				return nil, err
			}
			serialized = append(serialized, struct {
				Type  types.ViewType `json:"type"`
				State []byte         `json:"state"`
			}{
				Type:  vs.ViewType(),
				State: stateData,
			})
		}
	}
	return json.Marshal(serialized)
}

// DeserializeStack restores stack from JSON
func (ns *NavigationStack) DeserializeStack(data []byte) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	var serialized []struct {
		Type  types.ViewType `json:"type"`
		State []byte         `json:"state"`
	}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return err
	}
	newStack := make([]interface{}, 0, len(serialized))
	for _, item := range serialized {
		var state types.ViewState
		switch item.Type {
		case types.MenuStateType:
			state = &types.MenuViewState{}
		case types.ChatStateType:
			state = &types.ChatViewState{}
		case types.ModalStateType:
			state = &types.ModalViewState{}
		default:
			continue // Skip unknown types
		}
		if err := state.UnmarshalState(item.State); err != nil {
			return err
		}
		newStack = append(newStack, state)
	}
	if len(newStack) > 0 {
		if mm, ok := newStack[0].(types.ViewState); ok && mm.IsMainMenu() {
			ns.stack = newStack
		}
	}
	return nil
}

// CanPop returns true if the stack has more than one item (i.e., can pop without removing the root)
func (ns *NavigationStack) CanPop() bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return len(ns.stack) > 1
}

// Current returns the current (top) ViewState on the stack
func (ns *NavigationStack) Current() interface{} {
	return ns.Top()
}

// ShowModal is a stub for modal management (no-op)
func (ns *NavigationStack) ShowModal(modalType string, data interface{}) {
	// No-op for now
}

// HideModal is a stub for modal management (no-op)
func (ns *NavigationStack) HideModal() {
	// No-op for now
}

