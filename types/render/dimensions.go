// dimensions.go
// Common dimension and render strategy types for TUI components (menus, notices, inputs, etc.)
// Designed for use with Bubble Tea and lipgloss, supporting theming and responsive layouts.

package render

// Dimension defines width and height for a component
// Use -1 for flexible (auto) sizing
//
type Dimension struct {
	Width  int // in cells/columns, -1 for auto
	Height int // in rows, -1 for auto
}

// RenderStrategy defines layout and style hints for a component
//
type RenderStrategy struct {
	Dimension Dimension // base dimensions
	Padding   int       // uniform padding
	Margin    int       // uniform margin
	Border    bool      // whether to show border
	Centered  bool      // whether to center content
	ThemeKey  string    // key for theme/style lookup (e.g., "menu", "notice")
}

// Predefined strategies for common UI elements
var (
	MenuStrategy = RenderStrategy{
		Dimension: Dimension{Width: 40, Height: -1},
		Padding:   1,
		Margin:    1,
		Border:    true,
		Centered:  true,
		ThemeKey:  "menu",
	}
	NoticeStrategy = RenderStrategy{
		Dimension: Dimension{Width: 50, Height: -1},
		Padding:   2,
		Margin:    1,
		Border:    true,
		Centered:  true,
		ThemeKey:  "notice",
	}
	InputStrategy = RenderStrategy{
		Dimension: Dimension{Width: 30, Height: 3},
		Padding:   1,
		Margin:    1,
		Border:    true,
		Centered:  false,
		ThemeKey:  "input",
	}
)

// RenderViewType defines the main types of renderable views in the app
// e.g., menu, prompts, notices, input, modal
//
type RenderViewType string

const (
	MenuViewType    RenderViewType = "menu"
	PromptsViewType RenderViewType = "prompts"
	NoticesViewType RenderViewType = "notices"
	InputViewType   RenderViewType = "input"
	ModalViewType   RenderViewType = "modal"
)

// DefaultRenderStrategies maps each RenderViewType to a default RenderStrategy
var DefaultRenderStrategies = map[RenderViewType]RenderStrategy{
	MenuViewType:    MenuStrategy,
	PromptsViewType: MenuStrategy, // Can customize if needed
	NoticesViewType: NoticeStrategy,
	InputViewType:   InputStrategy,
	ModalViewType:   NoticeStrategy, // Can customize if needed
}

// You can add more strategies as needed for other component types.
