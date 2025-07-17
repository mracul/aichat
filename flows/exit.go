package flows
package flows

import (
	"aichat/components/modals"
	"aichat/components/modals/dialogs"
	"aichat/interfaces"
)

// FlowExitMenu is the canonical exit flow for menus and global actions
func FlowExitMenu(nav interfaces.Controller) {
	modal := dialogs.NewConfirmationModal(
		"Are you sure you want to exit?",
		[]modals.ModalOption{
			{
				Label: "Yes",
				OnSelect: func() {
					// Send a quit message via the controller
					if sender, ok := nav.(interface{ QuitApp() }); ok {
						sender.QuitApp()
					}
				},
			},
			{
				Label: "No",
				OnSelect: func() {
					nav.HideModal()
				},
			},
		},
		func() { nav.HideModal() },
		modals.ModalRenderConfig{},
	)
	nav.ShowModal("confirmation", modal)
}

