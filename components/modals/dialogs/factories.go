package dialogs
package dialogs

import (
	"aichat/components/modals"
)

const DefaultMenuModalWidth = 300
const DefaultMenuModalHeight = 400
const DefaultListModalWidth = (DefaultMenuModalWidth * 8) / 10   // 80%
const DefaultListModalHeight = (DefaultMenuModalHeight * 8) / 10 // 80%
const DefaultConfirmationModalWidth = 300
const DefaultConfirmationModalHeight = 150

// Factory for ConfirmationModal
func NewConfirmationModalFactory(message string, options []modals.ModalOption, closeSelf modals.CloseSelfFunc, config modals.ModalRenderConfig) *ConfirmationModal {
	return NewConfirmationModal(message, options, closeSelf, config)
}

// Factory for ListModal
func NewListModalFactory(title string, options []string, onSelect func(int), closeSelf func(), config modals.ModalRenderConfig) *ListModal {
	return &ListModal{
		Title:         title,
		Options:       options,
		Selected:      0,
		OnSelect:      onSelect,
		CloseSelfFunc: closeSelf,
		Config:        config,
		RegionWidth:   DefaultListModalWidth,
		RegionHeight:  DefaultListModalHeight,
	}
}

// Factory for MenuModal
func NewMenuModalFactory(title string, options []string, onSelect func(int), closeSelf func(), config modals.ModalRenderConfig) *MenuModal {
	return &MenuModal{
		Title:        title,
		Options:      options,
		Selected:     0,
		OnSelect:     onSelect,
		CloseSelf:    closeSelf,
		Config:       config,
		RegionWidth:  DefaultMenuModalWidth,
		RegionHeight: DefaultMenuModalHeight,
	}
}

// Factory for HelpModal
func NewHelpModalFactory(content string, closeSelf func(), regionWidth, regionHeight int, config modals.ModalRenderConfig) *HelpModal {
	return &HelpModal{
		Content:           content,
		CloseSelf:         closeSelf,
		RegionWidth:       DefaultListModalWidth,
		RegionHeight:      DefaultListModalHeight,
		ModalRenderConfig: config,
	}
}

