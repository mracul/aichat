// modals.go - Modal type definitions for all modal/popup types in the application.
// Defines the Modal interface and all modal data structures for use with the navigation/controller system.

package types

// Modal is the interface for all modal types.
type Modal interface {
	Type() ModalType
	View() string
	Update(msg interface{}, ctx Context, nav Controller) (Modal, interface{})
	OnShow()
	OnHide()
	IsClosable() bool
	CloseSelf()
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
	OnSelect  []func(ctx Context, nav Controller)
	CloseSelf func()
}

// InputModal - Single-line or multi-line text input dialog.
type InputModal struct {
	Prompt    string
	Value     string
	OnSubmit  func(value string, ctx Context, nav Controller)
	OnCancel  func(ctx Context, nav Controller)
	CloseSelf func()
}

// NoticeModal - Information or notice dialog with optional title.
type NoticeModal struct {
	Title     string
	Message   string
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}

// ErrorModal - Specialized notice for error messages.
type ErrorModal struct {
	Title     string
	Message   string
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}

// SelectionListModal - List selection dialog for choosing from a list of items.
type SelectionListModal struct {
	Title     string
	Items     []string
	OnSelect  func(index int, ctx Context, nav Controller)
	OnCancel  func(ctx Context, nav Controller)
	CloseSelf func()
}

// EditorModal - Multi-line text editor modal.
type EditorModal struct {
	Title     string
	Content   string
	OnSubmit  func(content string, ctx Context, nav Controller)
	OnCancel  func(ctx Context, nav Controller)
	CloseSelf func()
}

// HelpModal - Help/controls cheat sheet.
type HelpModal struct {
	Content   string
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}

// AboutModal - About/info modal.
type AboutModal struct {
	Content   string
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}

// GoodbyeModal - Goodbye/exit modal.
type GoodbyeModal struct {
	Message   string
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}

// CustomModal - For future extensibility.
type CustomModal struct {
	TypeName  string
	Data      map[string]interface{}
	OnClose   func(ctx Context, nav Controller)
	CloseSelf func()
}
