package modals
package modals

import (
	"aichat/types/render"
)

// ModalRenderConfig holds all theming and layout info for modals.
type ModalRenderConfig struct {
	ThemeMap   render.ThemeMap
	Strategies map[string]render.RenderStrategy
}

// RenderContentWithStrategy renders content using the config and a strategy key.
func (c ModalRenderConfig) RenderContentWithStrategy(content string, strategyKey string) string {
	return RenderModalBox(content, c, strategyKey)
}

// RenderModalBox renders a modal box with content, using the given config and strategy key.
func RenderModalBox(content string, config ModalRenderConfig, strategyKey string) string {
	strategy := config.Strategies[strategyKey]
	theme := config.ThemeMap[strategy.ThemeKey]
	return render.ApplyStrategy(content, strategy, theme)
}

// Additional helpers (e.g., RenderModalHeader, RenderModalFooter) can be added as needed.
