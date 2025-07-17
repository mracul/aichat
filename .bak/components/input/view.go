package input
package input

import (
	"aichat/models"
	"aichat/types"
	"fmt"
)

// RenderInputView renders the input area and its controls.
func RenderInputView(m *models.InputModel, ctx types.Context) string {
	inputBox := m.View()
	controls := m.GetControlSet(ctx)
	var controlHints string
	for _, ctrl := range controls.Controls {
		if ctrl.Key != 0 {
			controlHints += fmt.Sprintf("[%s: %s] ", ctrl.Name, ctrl.Key.String())
		} else {
			controlHints += fmt.Sprintf("[%s] ", ctrl.Name)
		}
	}
	return inputBox + "\n" + controlHints
}

// RenderViewWithControls renders any view and its controls if it implements ControlSetProvider.
func RenderViewWithControls(view types.Renderable, ctx types.Context) string {
	content := view.View()
	if provider, ok := view.(types.ControlSetProvider); ok {
		controls := provider.GetControlSet(ctx)
		var controlHints string
		for _, ctrl := range controls.Controls {
			if ctrl.Key != 0 {
				controlHints += fmt.Sprintf("[%s: %s] ", ctrl.Name, ctrl.Key.String())
			} else {
				controlHints += fmt.Sprintf("[%s] ", ctrl.Name)
			}
		}
		content += "\n" + controlHints
	}
	return content
}
