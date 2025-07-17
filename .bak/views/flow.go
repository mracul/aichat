package views
// flow.go - Modular, extensible multi-step modal flow system for TUI navigation
// Implements FlowViewState, ModalSet, FlowModal, FlowRunner, and robust handler integration
// Integrates with navigation/controller and menu system

package views

import (
	"aichat/types"

	tea "github.com/charmbracelet/bubbletea"
)

// FlowRunner provides control methods for modals to advance or cancel the flow
// Modals receive this via SetFlowRunner
// This interface is injected into each modal (or ModalSet)
type FlowRunner interface {
	Next(data map[string]any)
	Cancel()
}

// FlowModal is a modal that can participate in a flow
// It must implement the required methods for modal flows
// UI-only, signals via FlowRunner
// Unaware of navigation or flow internals
type FlowModal interface {
	Type() types.ViewType
	Init() tea.Cmd
	View() string
	SetFlowRunner(runner FlowRunner)
	UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd)
}

// ModalSet represents a logical grouping or single unit of modal steps
// Can be a single modal or a group/branch
// Manages rendering and input for one or more modals in a set
// Interacts with the Flow via ModalRunner
// Abstracts grouping/branching behaviors within a flow
type ModalSet interface {
	Init() tea.Cmd
	View() string
	UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (ModalSet, tea.Cmd)
	SetFlowRunner(runner FlowRunner)
	IsComplete() bool
}

// FlowViewState orchestrates progression through a sequence of ModalSet[]
// Accumulates and manages shared state (map[string]any)
// Controls flow lifecycle (begin, advance, cancel, complete)
// Implements ViewState interface for navigation stack integration
// Injects context and navigation controller dependencies
// Triggers onSuccess or onExit handlers with amassed data
// Handles observer registration/unregistration for ModalSets/FlowModals

type FlowViewState struct {
	Steps     []ModalSet
	StepIndex int
	State     map[string]any
	OnSuccess func(state map[string]any, ctx types.Context, nav types.Controller)
	OnExit    func(state map[string]any, ctx types.Context, nav types.Controller)
	ctx       types.Context
	nav       types.Controller
	Subject   types.Subject // The model/state to observe (e.g., chat, flow, etc.)
}

// NewFlowViewState creates a new flow with the given steps and handlers
func NewFlowViewState(
	steps []ModalSet,
	onSuccess func(state map[string]any, ctx types.Context, nav types.Controller),
	onExit func(state map[string]any, ctx types.Context, nav types.Controller),
	ctx types.Context,
	nav types.Controller,
	subject types.Subject, // new parameter for observer management
) *FlowViewState {
	f := &FlowViewState{
		Steps:     steps,
		StepIndex: 0,
		State:     make(map[string]any),
		OnSuccess: onSuccess,
		OnExit:    onExit,
		ctx:       ctx,
		nav:       nav,
		Subject:   subject,
	}
	if len(steps) > 0 {
		steps[0].SetFlowRunner(f)
		// Register as observer if applicable
		if subject != nil {
			if observer, ok := steps[0].(types.Observer); ok {
				subject.RegisterObserver(observer)
			}
		}
	}
	return f
}

// Type returns the ViewType for the flow (could define a new type if desired)
func (f *FlowViewState) Type() types.ViewType { return types.ModalStateType }

// IsMainMenu returns false for FlowViewState (not a main menu)
func (f *FlowViewState) IsMainMenu() bool { return false }

// MarshalState is a stub for state serialization
func (f *FlowViewState) MarshalState() ([]byte, error) { return nil, nil }

// UnmarshalState is a stub for state deserialization
func (f *FlowViewState) UnmarshalState([]byte) error { return nil }

// ViewType returns the ViewType for the flow (for interface compliance)
func (f *FlowViewState) ViewType() types.ViewType { return types.ModalStateType }

// Init is a no-op for the flow
func (f *FlowViewState) Init() tea.Cmd { return nil }

// View renders only the current ModalSet, centered on a blank background
func (f *FlowViewState) View() string {
	if f.StepIndex < len(f.Steps) {
		return f.Steps[f.StepIndex].View() // Assume modal is already centered
	}
	return ""
}

// Bubble Tea-compatible Update delegates to OOP-style UpdateWithContext
func (f *FlowViewState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f.UpdateWithContext(msg, f.ctx, f.nav)
}

// In UpdateWithContext, handle observer registration/unregistration on step change
func (f *FlowViewState) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	if f.StepIndex < len(f.Steps) {
		set, cmd := f.Steps[f.StepIndex].UpdateWithContext(msg, ctx, nav)
		f.Steps[f.StepIndex] = set
		if set.IsComplete() {
			// Unregister old step as observer if applicable
			if f.Subject != nil {
				if observer, ok := set.(types.Observer); ok {
					f.Subject.UnregisterObserver(observer)
				}
			}
			// Advance to next step
			f.StepIndex++
			if f.StepIndex < len(f.Steps) {
				f.Steps[f.StepIndex].SetFlowRunner(f)
				// Register new step as observer if applicable
				if f.Subject != nil {
					if observer, ok := f.Steps[f.StepIndex].(types.Observer); ok {
						f.Subject.RegisterObserver(observer)
					}
				}
			} else {
				// Flow complete
				if f.OnSuccess != nil {
					f.OnSuccess(f.State, f.ctx, f.nav)
				}
				f.nav.Pop() // Remove flow from stack
			}
		}
		return f, cmd
	}
	return f, nil
}

// FlowRunner implementation
func (f *FlowViewState) Next(data map[string]any) {
	// Merge new data into a fresh copy of state
	newState := make(map[string]any, len(f.State)+len(data))
	for k, v := range f.State {
		newState[k] = v
	}
	for k, v := range data {
		newState[k] = v
	}
	f.State = newState
	// Mark current ModalSet as complete (should be handled by ModalSet logic)
	// Step will advance on next Update
}

func (f *FlowViewState) Cancel() {
	if f.OnExit != nil {
		f.OnExit(f.State, f.ctx, f.nav)
	}
	f.nav.Pop() // Remove flow from stack
}

// =========================
// Example ModalSet and FlowModal Implementations
// =========================

// BasicModalSet wraps a single FlowModal and manages its lifecycle.
// It is complete when the modal signals Next or Cancel via the FlowRunner interface.
// This is the simplest ModalSet; more complex sets can group or branch multiple modals.
type BasicModalSet struct {
	Modal      FlowModal
	complete   bool
	flowRunner FlowRunner
}

// NewBasicModalSet creates a ModalSet for a single FlowModal
func NewBasicModalSet(modal FlowModal) *BasicModalSet {
	return &BasicModalSet{Modal: modal, complete: false}
}

func (s *BasicModalSet) Init() tea.Cmd {
	return s.Modal.Init()
}

func (s *BasicModalSet) View() string {
	return s.Modal.View()
}

func (s *BasicModalSet) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (ModalSet, tea.Cmd) {
	if s.complete {
		return s, nil
	}
	model, cmd := s.Modal.UpdateWithContext(msg, ctx, nav)
	if updated, ok := model.(FlowModal); ok {
		s.Modal = updated
	}
	return s, cmd
}

func (s *BasicModalSet) SetFlowRunner(runner FlowRunner) {
	s.flowRunner = runner
	s.Modal.SetFlowRunner(&modalSetRunner{s})
}

func (s *BasicModalSet) IsComplete() bool {
	return s.complete
}

// modalSetRunner is a FlowRunner implementation that marks the ModalSet as complete and forwards to the parent FlowRunner
// Used internally to decouple Modal from ModalSet and Flow
// This allows Modal to signal Next/Cancel without knowing about the flow
// (see sequence diagram: Modal->>FlowViewState: FlowRunner.Next(data))
type modalSetRunner struct {
	set *BasicModalSet
}

func (r *modalSetRunner) Next(data map[string]any) {
	if r.set.flowRunner != nil {
		r.set.complete = true
		r.set.flowRunner.Next(data)
	}
}

func (r *modalSetRunner) Cancel() {
	if r.set.flowRunner != nil {
		r.set.complete = true
		r.set.flowRunner.Cancel()
	}
}

// =========================
// Example FlowModal: InputFlowModal
// =========================

// InputFlowModal collects a string value for a given key
// Presents a prompt, collects user input, and signals Next or Cancel via FlowRunner
// Decoupled from navigation/flow internals; UI-only
// (see sequence diagram: User->>Modal: Input data or Cancel)
type InputFlowModal struct {
	Prompt     string
	Key        string
	Value      string
	flowRunner FlowRunner
	completed  bool
}

func NewInputFlowModal(prompt, key string) *InputFlowModal {
	return &InputFlowModal{Prompt: prompt, Key: key}
}

func (m *InputFlowModal) SetFlowRunner(runner FlowRunner) { m.flowRunner = runner }
func (m *InputFlowModal) Type() types.ViewType            { return types.ModalStateType }
func (m *InputFlowModal) Init() tea.Cmd                   { return nil }
func (m *InputFlowModal) View() string                    { return m.Prompt + ": " + m.Value }

// Bubble Tea-compatible Update delegates to OOP-style UpdateWithContext
func (m *InputFlowModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.UpdateWithContext(msg, nil, nil)
}

// UpdateWithContext handles user input and signals Next or Cancel
func (m *InputFlowModal) UpdateWithContext(msg tea.Msg, ctx types.Context, nav types.Controller) (tea.Model, tea.Cmd) {
	if m.completed {
		return m, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			if m.flowRunner != nil {
				m.completed = true
				m.flowRunner.Next(map[string]any{m.Key: m.Value})
			}
			return m, nil
		case tea.KeyEsc:
			if m.flowRunner != nil {
				m.completed = true
				m.flowRunner.Cancel()
			}
			return m, nil
		case tea.KeyBackspace:
			if len(m.Value) > 0 {
				m.Value = m.Value[:len(m.Value)-1]
			}
			return m, nil
		default:
			if len(keyMsg.String()) == 1 {
				m.Value += keyMsg.String()
			}
			return m, nil
		}
	}
	return m, nil
}

// =========================
// Concrete Flow: Create New Chat Flow
// =========================

// CreateNewChatFlow composes a multi-step flow for creating a new chat.
// Steps: Enter chat title, enter prompt, enter model (all as InputFlowModal for now).
// Accumulates state and calls onSuccess with the collected data.
// Usage: Call from a menu action or controller to launch the flow.
func CreateNewChatFlow(
	ctx types.Context,
	nav types.Controller,
	onSuccess func(state map[string]any, ctx types.Context, nav types.Controller),
) {
	steps := []ModalSet{
		NewBasicModalSet(NewInputFlowModal("Enter chat title", "title")),
		NewBasicModalSet(NewInputFlowModal("Enter prompt", "prompt")),
		NewBasicModalSet(NewInputFlowModal("Enter model", "model")),
	}
	flow := NewFlowViewState(
		steps,
		onSuccess,
		func(state map[string]any, ctx types.Context, nav types.Controller) {
			// Default onExit: do nothing (could show a notice or return to menu)
		},
		ctx, nav,
		nil, // No subject for this flow
	)
	nav.Push(flow)
}

// =========================
// Concrete Flow: Create New Chat (Simple)
// =========================

// CreateNewChatSimpleFlow launches a flow for creating a new chat with only a title.
// Uses defaults for prompt and model (can be parameterized or hardcoded).
// Usage: Call from a menu action or controller to launch the flow.
func CreateNewChatSimpleFlow(
	ctx types.Context,
	nav types.Controller,
	defaultPrompt string,
	defaultModel string,
	onSuccess func(state map[string]any, ctx types.Context, nav types.Controller),
) {
	steps := []ModalSet{
		NewBasicModalSet(NewInputFlowModal("Enter chat title", "title")),
	}
	flow := NewFlowViewState(
		steps,
		func(state map[string]any, ctx types.Context, nav types.Controller) {
			// Add defaults to state before calling onSuccess
			state["prompt"] = defaultPrompt
			state["model"] = defaultModel
			onSuccess(state, ctx, nav)
		},
		func(state map[string]any, ctx types.Context, nav types.Controller) {
			// Default onExit: do nothing (could show a notice or return to menu)
		},
		ctx, nav,
		nil, // No subject for this flow
	)
	nav.Push(flow)
}

// =========================
// Usage Example
// =========================
// In a menu action or controller:
// CreateNewChatFlow(ctx, nav, func(state map[string]any, ctx types.Context, nav types.Controller) {
//     // Handle the new chat creation with collected state
//     title := state["title"].(string)
//     prompt := state["prompt"].(string)
//     model := state["model"].(string)
//     // ... create chat, update UI, etc. ...
// })
// In a menu action or controller:
// CreateNewChatSimpleFlow(ctx, nav, "Default prompt text", "gpt-3.5-turbo", func(state map[string]any, ctx types.Context, nav types.Controller) {
//     // Handle the new chat creation with collected state
//     title := state["title"].(string)
//     prompt := state["prompt"].(string) // will be default
//     model := state["model"].(string)   // will be default
//     // ... create chat, update UI, etc. ...
// })
