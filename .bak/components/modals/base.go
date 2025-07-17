package modals
// base.go - Contains the BaseModal struct and ModalOption/CloseSelf types for reusable modal dialogs in the Bubble Tea UI.

package modals

// ModalOption represents a selectable option in a modal dialog.
type ModalOption struct {
	Label    string
	OnSelect func()
}

// CloseSelfFunc is a callback to restore previous app state when a modal is closed.
type CloseSelfFunc func()

// BaseModal is the base struct for all modal dialogs, containing message, options, close logic, and selection state.
type BaseModal struct {
	ModalRenderConfig
	Message      string
	Options      []ModalOption
	CloseSelf    CloseSelfFunc
	Selected     int
	RegionWidth  int
	RegionHeight int
}

// Init initializes the modal (Bubble Tea compatibility).
func (m *BaseModal) Init() {}

// Update handles Bubble Tea messages for the modal.
func (m *BaseModal) Update(msg interface{}) {}

// ViewRegion renders the modal content using the shared modal rendering utility.
func (m *BaseModal) ViewRegion(content string, strategyKey string) string {
	return RenderModalBox(content, m.ModalRenderConfig, strategyKey)
}

// RenderContentWithStrategy renders content using the given strategy key.
func (m *BaseModal) RenderContentWithStrategy(content string, strategyKey string) string {
	return RenderModalBox(content, m.ModalRenderConfig, strategyKey)
}
