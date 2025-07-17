package menus
package menus

import (
	"aichat/types"
	"aichat/interfaces"
	"aichat/types"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type ThemeMenuViewState struct {
	entries           types.MenuEntrySet
	previewIndex      int
	originalThemeName string
	currentThemeName  string
	themes            []map[string]interface{} // raw theme objects for flexibility
	themeNames        []string
	ctx               interfaces.Context
	nav               interfaces.Controller
	WindowWidth       int
	WindowHeight      int
}

// Add Generate Theme entry to the Themes menu
func NewThemeMenuViewState(ctx interfaces.Context, nav interfaces.Controller, width, height int) *ThemeMenuViewState {
	themes, names := loadAllThemes()
	currentTheme, _ := app.GetCurrentTheme()
	entries := buildThemeEntries(names)
	// Add Generate Theme entry at the top
	entries = append([]types.MenuEntry{{
		Text:        "Generate Theme",
		Description: "Create a new theme from two colors",
		Action: func(ctx interfaces.Context, nav interfaces.Controller) error {
			// Start the generate theme flow
			nav.Push(NewGenerateThemeFlowState(ctx, nav, width, height))
			return nil
		},
	}}, entries...)
	return &ThemeMenuViewState{
		entries:           entries,
		previewIndex:      0,
		originalThemeName: currentTheme,
		currentThemeName:  currentTheme,
		themes:            themes,
		themeNames:        append([]string{"Generate Theme"}, names...),
		ctx:               ctx,
		nav:               nav,
		WindowWidth:       width,
		WindowHeight:      height,
	}
}

func buildThemeEntries(names []string) types.MenuEntrySet {
	entries := make(types.MenuEntrySet, len(names)+1)
	for i, name := range names {
		entries[i] = types.MenuEntry{
			Text:        name,
			Description: "Preview and set this theme",
		}
	}
	entries[len(names)] = types.MenuEntry{
		Text:        "Back",
		Description: "Return to settings",
		Next:        types.SettingsMenu,
	}
	return entries
}

func loadAllThemes() ([]map[string]interface{}, []string) {
	data, err := ioutil.ReadFile(".config/themes.json")
	if err != nil {
		return nil, nil
	}
	var themes []map[string]interface{}
	_ = json.Unmarshal(data, &themes)
	names := make([]string, len(themes))
	for i, t := range themes {
		if n, ok := t["name"].(string); ok {
			names[i] = n
		}
	}
	return themes, names
}

func (t *ThemeMenuViewState) Type() types.ViewType          { return types.MenuStateType }
func (t *ThemeMenuViewState) IsMainMenu() bool              { return false }
func (t *ThemeMenuViewState) ViewType() types.ViewType      { return types.MenuStateType }
func (t *ThemeMenuViewState) MarshalState() ([]byte, error) { return nil, nil }
func (t *ThemeMenuViewState) UnmarshalState([]byte) error   { return nil }
func (t *ThemeMenuViewState) Init() tea.Cmd                 { return nil }

func (t *ThemeMenuViewState) View() string {
	// Use the preview theme for styling
	// (Integration with the app's theme system is required here)
	out := "\nSelect a theme (Up/Down to preview, Enter to set, Esc to cancel):\n\n"
	for i, name := range t.themeNames {
		prefix := "  "
		if i == t.previewIndex {
			prefix = "> "
		}
		out += prefix + name + "\n"
	}
	out += "\n[Up/Down] Preview  [Enter] Set  [Esc] Cancel"
	return out
}

func (t *ThemeMenuViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if t.previewIndex > 0 {
				t.previewIndex--
			} else {
				t.previewIndex = len(t.themeNames) - 1
			}
			t.applyPreviewTheme()
		case "down", "j":
			if t.previewIndex < len(t.themeNames)-1 {
				t.previewIndex++
			} else {
				t.previewIndex = 0
			}
			t.applyPreviewTheme()
		case "enter":
			selected := t.themeNames[t.previewIndex]
			app.SetCurrentTheme(selected)
			// Show notice modal (integration required)
			// nav.ShowModal("notice", "Theme set to "+selected)
			return t, nil
		case "esc":
			t.restoreOriginalTheme()
			// nav.Pop() or return to settings
			return t, nil
		}
	}
	return t, nil
}

func (t *ThemeMenuViewState) applyPreviewTheme() {
	selected := t.themeNames[t.previewIndex]
	// Integration: update app's theme to preview this theme
	app.SetCurrentTheme(selected) // For now, just set as current (should be preview only)
}

func (t *ThemeMenuViewState) restoreOriginalTheme() {
	app.SetCurrentTheme(t.originalThemeName)
}

// --- Generate Theme Flow ---

// --- Generate Theme Flow State ---
type GenerateThemeFlowState struct {
	step         int // 0: color1, 1: color2, 2: waiting, 3: preview, 4: name
	colors       [2]string
	input        string
	errorMsg     string
	previewTheme map[string]interface{}
	ctx          interfaces.Context
	nav          interfaces.Controller
	WindowWidth  int
	WindowHeight int
}

func NewGenerateThemeFlowState(ctx interfaces.Context, nav interfaces.Controller, width, height int) *GenerateThemeFlowState {
	return &GenerateThemeFlowState{
		step:         0,
		ctx:          ctx,
		nav:          nav,
		WindowWidth:  width,
		WindowHeight: height,
	}
}

func (g *GenerateThemeFlowState) Type() types.ViewType          { return types.MenuStateType }
func (g *GenerateThemeFlowState) IsMainMenu() bool              { return false }
func (g *GenerateThemeFlowState) ViewType() types.ViewType      { return types.MenuStateType }
func (g *GenerateThemeFlowState) MarshalState() ([]byte, error) { return nil, nil }
func (g *GenerateThemeFlowState) UnmarshalState([]byte) error   { return nil }
func (g *GenerateThemeFlowState) Init() tea.Cmd                 { return nil }

func (g *GenerateThemeFlowState) View() string {
	switch g.step {
	case 0:
		return "Enter primary accent color (hex, e.g. #3498db):\n" + g.input + "\n"
	case 1:
		return "Enter secondary accent color (hex, e.g. #2ecc71):\n" + g.input + "\n"
	case 2:
		return "Generating theme...\n"
	case 3:
		// Show preview (simplified)
		preview, _ := json.MarshalIndent(g.previewTheme, "", "  ")
		return "Preview generated theme:\n" + string(preview) + "\n\nPress Enter to name and save, Esc to cancel."
	case 4:
		return "Enter a name for your theme:\n" + g.input + "\n"
	case 5:
		return "Theme saved! Returning to themes menu..."
	}
	if g.errorMsg != "" {
		return "Error: " + g.errorMsg + "\nPress any key to return to themes menu."
	}
	return ""
}

func (g *GenerateThemeFlowState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		if g.errorMsg != "" {
			g.nav.Pop()
			return g, nil
		}
		switch g.step {
		case 0, 1:
			if m.Type == tea.KeyEnter {
				if !isValidHexColor(g.input) {
					g.errorMsg = "Invalid color code. Please enter a valid hex color (e.g. #3498db)."
					return g, nil
				}
				g.colors[g.step] = g.input
				g.input = ""
				g.step++
				return g, nil
			} else if m.Type == tea.KeyBackspace && len(g.input) > 0 {
				g.input = g.input[:len(g.input)-1]
				return g, nil
			} else if len(m.String()) == 1 || (len(m.String()) == 2 && m.String()[0] == '#') {
				g.input += m.String()
				return g, nil
			}
		case 2:
			// Ignore input while waiting
			return g, nil
		case 3:
			if m.Type == tea.KeyEnter {
				g.step = 4
				g.input = ""
				return g, nil
			} else if m.Type == tea.KeyEsc {
				g.nav.Pop()
				return g, nil
			}
		case 4:
			if m.Type == tea.KeyEnter && g.input != "" {
				g.saveTheme()
				g.step = 5
				return g, nil
			} else if m.Type == tea.KeyBackspace && len(g.input) > 0 {
				g.input = g.input[:len(g.input)-1]
				return g, nil
			} else if len(m.String()) == 1 {
				g.input += m.String()
				return g, nil
			}
		case 5:
			g.nav.Pop()
			return g, nil
		}
	}
	// After color input, trigger API call
	if g.step == 2 {
		g.generateThemeFromAPI()
	}
	return g, nil
}

func (g *GenerateThemeFlowState) UpdateWithContext(msg tea.Msg, ctx interfaces.Context, nav interfaces.Controller) (tea.Model, tea.Cmd) {
	return g.Update(msg)
}

func isValidHexColor(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')) {
			return false
		}
	}
	return true
}

func (g *GenerateThemeFlowState) generateThemeFromAPI() {
	g.step = 2
	prompt := `You are a color scheme designer.  
Given two input colors, generate a complete, aesthetic, and readable color scheme for a terminal/chat application.  
The output must be a single JSON object matching the schema below, with all color values as hex codes.  
The scheme should be visually consistent, accessible, and harmonious, using the two input colors as the primary accents.  
Use the first color as the main accent (for highlights, focused borders, etc.) and the second as a secondary accent (for backgrounds, window borders, etc.).  
Ensure good contrast and readability for all text and backgrounds.

Special instructions:  
- The error and notice modals should be appropriately colored to match their desired purpose (e.g., error modals should use red or a high-contrast warning color, notice modals should use yellow/orange or another attention-grabbing color).  
- You may use common colors for these modals, or derive suitable colors through contrast or other methods you see fit, to ensure their intent is clear and visually distinct.
- Take inspiration from popular nvim text editor themes such as Monokai, Dracula, Catppuccin, Gruvbox, Solarized, and Tokyo Night. The overall palette should feel modern and in line with these themes.

Schema:
{"name": "Custom Theme", "modal": {"error": {"textColor": "", "highlightTextColor": "", "windowBorder": ""}, "input": {"textColor": "", "highlightTextColor": "", "windowBorder": ""}, "notice": {"textColor": "", "highlightTextColor": "", "windowBorder": ""}}, "menu": {"textColor": "", "highlightTextColor": "", "windowBorder": ""}, "window": {"focusedBorder": "", "unfocusedBorder": ""}, "appTextColor": ""}

Instructions:  
- Fill in all fields with appropriate hex color codes.  
- Use the two input colors as the basis for the palette, but you may derive tints/shades as needed for contrast and harmony.  
- Ensure all text is readable against its background.  
- The theme should look modern and visually pleasing.  
- Output only the JSON object, no explanations.

Example input:  
Color 1: ` + g.colors[0] + `  
Color 2: ` + g.colors[1] + `
`
	// TODO: Get active key, endpoint, and default model
	apiKey := os.Getenv("API_KEY")        // Replace with actual key retrieval
	endpoint := os.Getenv("API_ENDPOINT") // Replace with actual endpoint retrieval
	model := "gpt-3.5-turbo"              // Replace with actual model retrieval
	if apiKey == "" || endpoint == "" || model == "" {
		g.errorMsg = "Missing API key, endpoint, or model."
		g.step = -1
		return
	}
	// Prepare request
	body := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "system", "content": prompt}},
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		g.errorMsg = "Failed to create API request."
		g.step = -1
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.errorMsg = "API request failed."
		g.step = -1
		return
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	// Parse response (assume response is just the JSON theme)
	var theme map[string]interface{}
	if err := json.Unmarshal(respBody, &theme); err != nil {
		g.errorMsg = "Failed to parse theme from API response."
		g.step = -1
		return
	}
	g.previewTheme = theme
	g.step = 3
}

func (g *GenerateThemeFlowState) saveTheme() {
	// Load existing themes
	data, err := ioutil.ReadFile(".config/themes.json")
	var themes []map[string]interface{}
	if err == nil {
		_ = json.Unmarshal(data, &themes)
	}
	// Set the name
	g.previewTheme["name"] = g.input
	// Append and save
	themes = append(themes, g.previewTheme)
	newData, _ := json.MarshalIndent(themes, "", "  ")
	_ = ioutil.WriteFile(".config/themes.json", newData, 0644)
}

