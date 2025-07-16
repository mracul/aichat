package navigation

import (
	"aichat/types"
	"log/slog"
	"sync"
)

type ControllerImpl struct {
	stack  []types.ViewState
	modals []types.ViewState // Modal stack, can be refined
	logger *slog.Logger
	mu     sync.RWMutex
}

func NewController(initialView types.ViewState, logger *slog.Logger) *ControllerImpl {
	return &ControllerImpl{
		stack:  []types.ViewState{initialView},
		modals: []types.ViewState{},
		logger: logger,
	}
}

func (nc *ControllerImpl) Push(view types.ViewState) {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	nc.stack = append(nc.stack, view)
	if nc.logger != nil {
		nc.logger.Debug("Navigation push", "view", view.ViewType(), "stack_depth", len(nc.stack))
	}
}

func (nc *ControllerImpl) Pop() types.ViewState {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.stack) <= 1 {
		if nc.logger != nil {
			nc.logger.Warn("Attempted to pop root view")
		}
		return nc.stack[0]
	}
	nc.stack = nc.stack[:len(nc.stack)-1]
	current := nc.stack[len(nc.stack)-1]
	if nc.logger != nil {
		nc.logger.Debug("Navigation pop", "view", current.ViewType(), "stack_depth", len(nc.stack))
	}
	return current
}

func (nc *ControllerImpl) Replace(view types.ViewState) {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.stack) == 0 {
		nc.stack = []types.ViewState{view}
	} else {
		nc.stack[len(nc.stack)-1] = view
	}
	if nc.logger != nil {
		nc.logger.Debug("Navigation replace", "view", view.ViewType(), "stack_depth", len(nc.stack))
	}
}

func (nc *ControllerImpl) ShowModal(modalType types.ModalType, data interface{}) {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	// For now, just log; actual modal management to be implemented
	if nc.logger != nil {
		nc.logger.Info("Show modal", "modalType", modalType)
	}
	// TODO: Push modal view state to modals stack
}

func (nc *ControllerImpl) HideModal() {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.modals) > 0 {
		nc.modals = nc.modals[:len(nc.modals)-1]
		if nc.logger != nil {
			nc.logger.Info("Hide modal", "remaining_modals", len(nc.modals))
		}
	}
}

func (nc *ControllerImpl) Current() types.ViewState {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.stack[len(nc.stack)-1]
}

func (nc *ControllerImpl) CanPop() bool {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return len(nc.stack) > 1
}

func (nc *ControllerImpl) Resize(w, h int) {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if len(nc.stack) > 0 {
		if resizer, ok := nc.stack[len(nc.stack)-1].(interface{ Resize(int, int) }); ok {
			resizer.Resize(w, h)
		}
	}
}
