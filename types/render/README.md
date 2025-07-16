# Render Types and Strategies

This package provides reusable types and utilities for rendering TUI components with consistent dimensions, layout, and theming using Bubble Tea and lipgloss.

## Files
- `dimensions.go`: Defines `Dimension` and `RenderStrategy` types, with presets for common UI elements (menu, notice, input, etc.).
- `theme.go`: Defines the `Theme` type and `ThemeMap` for color theming.
- `strategy.go`: Provides the `Renderable` interface and `ApplyStrategy` utility for rendering content with a strategy and theme.

## Usage
- Assign a `RenderStrategy` (e.g., `MenuStrategy`) and a `Theme` to your view or component.
- Use `ApplyStrategy(content, strategy, theme)` to render the view with the specified layout and colors.
- Implement the `Renderable` interface for custom views that need to support strategies.

## Best Practices
- Keep all dimension and style logic in these types for consistency.
- Use the `ThemeKey` in `RenderStrategy` to look up the correct theme for each component.
- Extend with new strategies and themes as your UI grows. 