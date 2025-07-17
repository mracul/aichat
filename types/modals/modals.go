// modals.go - Modal type definitions for all modal/popup types in the application.
// Defines the Modal interface and all modal data structures for use with the navigation/controller system.

package types

import (
	"aichat/interfaces"
)

// ModalLifecycle defines show/hide hooks
type ModalLifecycle interface {
	OnShow()
	OnHide()
}

// ModalRenderable defines rendering
type ModalRenderable interface {
	View() string
}

// ModalUpdatable defines update logic
type ModalUpdatable interface {
	Update(msg interface{}, ctx Context, nav interfaces.Controller) (Modal, interface{})
}

// Closable defines close logic
type Closable interface {
	IsClosable() bool
	CloseSelf()
}

// Modal composes all the above
type Modal interface {
	ModalLifecycle
	ModalRenderable
	ModalUpdatable
	Closable
	Type() ModalType
}

// ModalType constants for modal type identification.
const (
	ModalTypeConfirmation ModalType = "confirmation"
	ModalTypeInput        ModalType = "input"
	ModalTypeNotice       ModalType = "notice"
	ModalTypeError        ModalType = "error"
	ModalTypeSelection    ModalType = "selection"
	ModalTypeHelp         ModalType = "help"
	ModalTypeAbout        ModalType = "about"
	ModalTypeGoodbye      ModalType = "goodbye"
	ModalTypeEditor       ModalType = "editor"
	ModalTypeCustom       ModalType = "custom"
)

// ConfirmationModal - Yes/No or multi-choice confirmation dialog.
type ConfirmationModal struct {
	Message   string
	Choices   []string
	OnSelect  []func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
	Selected  int // Default selected option (0=first, 1=second, etc.)
}

// InputModal - Single-line or multi-line text input dialog.
type InputModal struct {
	Prompt    string
	Value     string
	OnSubmit  func(value string, ctx Context, nav interfaces.Controller)
	OnCancel  func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// NoticeModal - Information or notice dialog with optional title.
type NoticeModal struct {
	Title     string
	Message   string
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// ErrorModal - Specialized notice for error messages.
type ErrorModal struct {
	Title     string
	Message   string
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// SelectionListModal - List selection dialog for choosing from a list of items.
type SelectionListModal struct {
	Title     string
	Items     []string
	OnSelect  func(index int, ctx Context, nav interfaces.Controller)
	OnCancel  func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// EditorModal - Multi-line text editor modal.
type EditorModal struct {
	Title     string
	Content   string
	OnSubmit  func(content string, ctx Context, nav interfaces.Controller)
	OnCancel  func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// HelpModal - Help/controls cheat sheet.
type HelpModal struct {
	Content   string
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// AboutModal - About/info modal.
type AboutModal struct {
	Content   string
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// GoodbyeModal - Goodbye/exit modal.
type GoodbyeModal struct {
	Message   string
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}

// CustomModal - For future extensibility.
type CustomModal struct {
	TypeName  string
	Data      map[string]interface{}
	OnClose   func(ctx Context, nav interfaces.Controller)
	CloseSelf func()
}
