package types

// ViewType is an enum for different view state types.
type ViewType int

const (
	MenuStateType ViewType = iota
	ChatStateType
	ModalStateType
)

// ViewState interface is defined in view_state.go
// This file contains other interfaces and types

// Controller interface for navigation stack and modal management
// (Moved from src/navigation/interfaces.go to break import cycles)
type Controller interface {
	Push(view ViewState)
	Pop() ViewState
	Replace(view ViewState)
	ShowModal(modalType ModalType, data interface{})
	HideModal()
	Current() ViewState
	CanPop() bool
}

type Context interface {
	App() interface{}     // *UnifiedAppModel, avoid import cycle
	GUI() interface{}     // *GUIAppModel, avoid import cycle
	Storage() interface{} // NavigationStorage, avoid import cycle
	Config() interface{}  // *AppConfig, avoid import cycle
	Logger() interface{}  // *slog.Logger, avoid import cycle
}
