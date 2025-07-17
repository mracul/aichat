package modals
// ai_key_test.go - Modal for testing AI API keys with spinner and result
// Returns a key:value pair on completion for flows to append to their data

package modals

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AIKeyTestModal displays a spinner and result for API key testing.
// On completion, returns a key:value pair: {"ai_key_test_result": {success, code, msg}}
type AIKeyTestModal struct {
	KeyTitle   string
	Testing    bool
	Success    bool
	ErrorMsg   string
	ResultCode string
	SpinnerIdx int
	OnComplete func(result map[string]any)
	Timeout    time.Duration
	startTime  time.Time
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type spinnerTickMsg struct{}
type aiKeyTestResultMsg struct {
	Success   bool
	Code, Msg string
}

// NewAIKeyTestModal creates a new modal for testing an AI key
// OnComplete is called with a result map: {"ai_key_test_result": {success, code, msg}}
func NewAIKeyTestModal(keyTitle string, onComplete func(result map[string]any)) *AIKeyTestModal {
	return &AIKeyTestModal{
		KeyTitle:   keyTitle,
		Testing:    true,
		OnComplete: onComplete,
		Timeout:    10 * time.Second,
		startTime:  time.Now(),
	}
}

func (m *AIKeyTestModal) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (m *AIKeyTestModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerTickMsg:
		if m.Testing {
			m.SpinnerIdx = (m.SpinnerIdx + 1) % len(spinnerFrames)
			if time.Since(m.startTime) > m.Timeout {
				m.Testing = false
				m.Success = false
				m.ErrorMsg = "Timeout"
				m.ResultCode = "timeout"
				if m.OnComplete != nil {
					m.OnComplete(map[string]any{"ai_key_test_result": map[string]any{"success": false, "code": "timeout", "msg": "API key test timed out"}})
				}
			}
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return spinnerTickMsg{}
			})
		}
	case aiKeyTestResultMsg:
		m.Testing = false
		m.Success = msg.Success
		m.ErrorMsg = msg.Msg
		m.ResultCode = msg.Code
		if m.OnComplete != nil {
			m.OnComplete(map[string]any{"ai_key_test_result": map[string]any{"success": msg.Success, "code": msg.Code, "msg": msg.Msg}})
		}
	}
	return m, nil
}

func (m *AIKeyTestModal) View() string {
	if m.Testing {
		return fmt.Sprintf("Testing %s %s", m.KeyTitle, spinnerFrames[m.SpinnerIdx])
	}
	if m.Success {
		return fmt.Sprintf("✅ Success: %s", m.ErrorMsg)
	}
	return fmt.Sprintf("❌ Error: %s", m.ErrorMsg)
}

// Usage:
// modal := NewAIKeyTestModal("My API Key", func(result map[string]any) {
//     // result["ai_key_test_result"] contains {success, code, msg}
//     // Flows can append this to their data on success
// })

