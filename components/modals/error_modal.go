package modals
// Package modals provides animated error handling components with dynamic feedback
package modals

import (
	"fmt"
	"strings"
	"time"

	"aichat/errors"

	tea "github.com/charmbracelet/bubbletea"
)

// =====================================================================================
// 🎭 Animated Error Modal - Dynamic Feedback with Live Updates
// =====================================================================================

type ErrorModalState string

const (
	ErrorStateInitial   ErrorModalState = "initial"
	ErrorStateAnalyzing ErrorModalState = "analyzing"
	ErrorStateRetrying  ErrorModalState = "retrying"
	ErrorStateResolved  ErrorModalState = "resolved"
	ErrorStateFailed    ErrorModalState = "failed"
)

type ErrorModal struct {
	Error       *errors.DomainError
	State       ErrorModalState
	Attempts    int
	MaxAttempts int
	RetryConfig errors.RetryConfig
	OnRetry     func() error
	OnResolve   func()
	OnDismiss   func()

	// Animation state
	SpinnerIndex int
	LastUpdate   time.Time
	Elapsed      time.Duration

	// UI dimensions
	Width  int
	Height int

	BaseModal
}

// Commands for animation and state updates
type errorTickMsg time.Time
type errorRetryMsg struct{ attempt int }
type errorResolveMsg struct{ success bool }

func errorTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return errorTickMsg(t)
	})
}

func errorRetry(attempt int) tea.Cmd {
	return func() tea.Msg {
		return errorRetryMsg{attempt: attempt}
	}
}

func errorResolve(success bool) tea.Cmd {
	return func() tea.Msg {
		return errorResolveMsg{success: success}
	}
}

// NewErrorModal creates a new animated error modal
func NewErrorModal(err *errors.DomainError, retryConfig errors.RetryConfig, onRetry func() error, config ModalRenderConfig) *ErrorModal {
	return &ErrorModal{
		Error:       err,
		State:       ErrorStateInitial,
		RetryConfig: retryConfig,
		OnRetry:     onRetry,
		MaxAttempts: retryConfig.MaxAttempts,
		Width:       60,
		Height:      15,
		LastUpdate:  time.Now(),
		BaseModal:   BaseModal{ModalRenderConfig: config, RegionWidth: 60, RegionHeight: 15},
	}
}

func (m *ErrorModal) Init() tea.Cmd {
	return tea.Batch(
		errorTick(),
		m.startAnalysis(),
	)
}

func (m *ErrorModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", "q":
			if m.State == ErrorStateFailed || m.State == ErrorStateResolved {
				if m.OnDismiss != nil {
					m.OnDismiss()
				}
				return m, tea.Quit
			}
		case "r":
			if m.State == ErrorStateFailed && m.Error.Retryable {
				return m, m.startRetry()
			}
		}

	case errorTickMsg:
		m.SpinnerIndex = (m.SpinnerIndex + 1) % 4
		m.Elapsed = time.Since(m.LastUpdate)
		return m, errorTick()

	case errorRetryMsg:
		m.Attempts = msg.attempt
		if msg.attempt <= m.MaxAttempts {
			return m, m.performRetry()
		} else {
			return m, errorResolve(false)
		}

	case errorResolveMsg:
		if msg.success {
			m.State = ErrorStateResolved
		} else {
			m.State = ErrorStateFailed
		}
		return m, nil
	}

	return m, nil
}

func (m *ErrorModal) View() string {
	var content strings.Builder

	// Header with error type and status
	header := m.renderHeader()
	content.WriteString(header)
	content.WriteString("\n\n")

	// Main content based on state
	switch m.State {
	case ErrorStateInitial:
		content.WriteString(m.renderInitial())
	case ErrorStateAnalyzing:
		content.WriteString(m.renderAnalyzing())
	case ErrorStateRetrying:
		content.WriteString(m.renderRetrying())
	case ErrorStateResolved:
		content.WriteString(m.renderResolved())
	case ErrorStateFailed:
		content.WriteString(m.renderFailed())
	}

	// Footer with controls
	footer := m.renderFooter()
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap in modal box using shared utility
	return m.RenderContentWithStrategy(content.String(), "modalBox")
}

// =====================================================================================
// 🎨 UI Rendering Methods
// =====================================================================================

func (m *ErrorModal) renderHeader() string {
	errorType := string(m.Error.Type)
	status := m.getStatusText()
	header := fmt.Sprintf("%s • %s", strings.ToUpper(errorType), status)
	return m.RenderContentWithStrategy(header, "header")
}

func (m *ErrorModal) renderInitial() string {
	content := m.Error.UserMsg + "\n\n"
	if len(m.Error.Details) > 0 {
		content += "Analyzing error details..."
	}
	return m.RenderContentWithStrategy(content, "error")
}

func (m *ErrorModal) renderAnalyzing() string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinners[m.SpinnerIndex%len(spinners)]
	content := fmt.Sprintf("%s Analyzing error...\n\n%s", spinner, m.renderErrorDetails())
	return m.RenderContentWithStrategy(content, "analyzing")
}

func (m *ErrorModal) renderRetrying() string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinners[m.SpinnerIndex%len(spinners)]
	content := fmt.Sprintf("%s Retrying... (Attempt %d/%d)\n\n%s", spinner, m.Attempts, m.MaxAttempts, m.renderRetryProgress())
	return m.RenderContentWithStrategy(content, "retrying")
}

func (m *ErrorModal) renderResolved() string {
	content := "✅ Issue resolved!\n\nThe problem has been fixed automatically."
	return m.RenderContentWithStrategy(content, "resolved")
}

func (m *ErrorModal) renderFailed() string {
	content := "❌ Unable to resolve automatically\n\n" + m.Error.UserMsg + "\n\n"
	if m.Error.Retryable {
		content += "Press 'r' to retry manually"
	}
	return m.RenderContentWithStrategy(content, "failed")
}

func (m *ErrorModal) renderFooter() string {
	var controls []string
	switch m.State {
	case ErrorStateFailed:
		if m.Error.Retryable {
			controls = append(controls, "r: Retry")
		}
		controls = append(controls, "Enter: Dismiss")
	case ErrorStateResolved:
		controls = append(controls, "Enter: Continue")
	default:
		controls = append(controls, "Esc: Cancel")
	}
	return m.RenderContentWithStrategy(strings.Join(controls, " • "), "footer")
}

func (m *ErrorModal) renderErrorDetails() string {
	if len(m.Error.Details) == 0 {
		return ""
	}
	var details string
	details += "Error Details:\n"
	for k, v := range m.Error.Details {
		details += fmt.Sprintf("  %s: %v\n", k, v)
	}
	return m.RenderContentWithStrategy(details, "details")
}

func (m *ErrorModal) renderRetryProgress() string {
	content := fmt.Sprintf("Retry Strategy: %s\nMax Attempts: %d\nInitial Delay: %v", m.RetryConfig.Backoff, m.RetryConfig.MaxAttempts, m.RetryConfig.InitialDelay)
	return m.RenderContentWithStrategy(content, "progress")
}

// =====================================================================================
// 🔄 State Management Methods
// =====================================================================================

func (m *ErrorModal) getStatusText() string {
	switch m.State {
	case ErrorStateInitial:
		return "Initializing"
	case ErrorStateAnalyzing:
		return "Analyzing"
	case ErrorStateRetrying:
		return fmt.Sprintf("Retrying (%d/%d)", m.Attempts, m.MaxAttempts)
	case ErrorStateResolved:
		return "Resolved"
	case ErrorStateFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

func (m *ErrorModal) startAnalysis() tea.Cmd {
	m.State = ErrorStateAnalyzing
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		if m.Error.Retryable {
			return errorRetry(1)
		} else {
			return errorResolve(false)
		}
	})
}

func (m *ErrorModal) startRetry() tea.Cmd {
	m.State = ErrorStateRetrying
	return errorRetry(1)
}

func (m *ErrorModal) performRetry() tea.Cmd {
	return func() tea.Msg {
		if m.OnRetry != nil {
			if err := m.OnRetry(); err == nil {
				return errorResolve(true)
			}
		}

		// Schedule next retry
		delay := m.calculateDelay(m.Attempts)
		time.Sleep(delay)

		return errorRetry(m.Attempts + 1)
	}
}

func (m *ErrorModal) calculateDelay(n int) time.Duration {
	switch m.RetryConfig.Backoff {
	case errors.LinearBackoff:
		return min(m.RetryConfig.InitialDelay*time.Duration(n), m.RetryConfig.MaxDelay)
	case errors.ExponentialBackoff:
		return min(m.RetryConfig.InitialDelay*time.Duration(1<<uint(n-1)), m.RetryConfig.MaxDelay)
	default:
		return m.RetryConfig.InitialDelay
	}
}

// =====================================================================================
// 🎯 Public Interface Methods
// =====================================================================================

// SetCallbacks sets the callback functions for the modal
func (m *ErrorModal) SetCallbacks(onResolve, onDismiss func()) {
	m.OnResolve = onResolve
	m.OnDismiss = onDismiss
}

// GetError returns the current error
func (m *ErrorModal) GetError() *errors.DomainError {
	return m.Error
}

// IsResolved returns true if the error was resolved
func (m *ErrorModal) IsResolved() bool {
	return m.State == ErrorStateResolved
}

// IsFailed returns true if the error handling failed
func (m *ErrorModal) IsFailed() bool {
	return m.State == ErrorStateFailed
}

