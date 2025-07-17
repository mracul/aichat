package models
package models

import (
	"strings"

	"aichat/interfaces"
)

type InputModel struct {
	Buffer    string
	Cursor    int
	Focused   bool
	Quitting  bool
	Width     int
	Height    int
	Message   string
	observers []interfaces.Observer // Observer pattern
}

// NewInputModel constructs a new InputModel.
func NewInputModel() *InputModel {
	return &InputModel{}
}

// Ensure InputModel implements types.Renderable and types.Updatable
var _ interface {
	interfaces.Renderable
	interfaces.Updatable
} = (*InputModel)(nil)

// Ensure InputModel implements types.ControlSetProvider
var _ interfaces.ControlSetProvider = (*InputModel)(nil)

func (m *InputModel) GetControlSet(ctx interfaces.Context) interfaces.ControlSet {
	if m.Focused {
		return interfaces.ControlSet{
			Controls: []interfaces.ControlType{
				{Name: "Send", Key: interfaces.KeyEnter, Action: nil},
				{Name: "Paste", Key: interfaces.KeyCtrlV, Action: nil},
				{Name: "Cut", Key: interfaces.KeyCtrlX, Action: nil},
				{Name: "Copy", Key: interfaces.KeyCtrlC, Action: nil},
				{Name: "Unfocus", Key: interfaces.KeyEsc, Action: nil},
			},
		}
	} else {
		return interfaces.ControlSet{
			Controls: []interfaces.ControlType{
				{Name: "Focus Input", Key: interfaces.KeyCtrlI, Action: nil},
			},
		}
	}
}

func (m *InputModel) Init() interfaces.Cmd { return nil }

func (m *InputModel) Update(msg interfaces.Msg) (interfaces.Model, interfaces.Cmd) {
	if !m.Focused {
		if keyMsg, ok := msg.(interfaces.KeyMsg); ok && keyMsg.String() == "ctrl+i" {
			m.Focused = true
		}
		return m, nil
	}
	if keyMsg, ok := msg.(interfaces.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+q", "esc":
			return m, interfaces.Quit
		case "ctrl+c":
			// Clipboard logic omitted for models package
			return m, nil
		case "ctrl+x":
			m.Buffer = ""
			m.Cursor = 0
			return m, nil
		case "ctrl+v":
			// Clipboard logic omitted for models package
			return m, nil
		case "ctrl+i":
			m.Focused = true
			return m, nil
		case "enter":
			if keyMsg.String() == "shift+enter" || keyMsg.String() == "alt+enter" {
				m.Buffer = m.Buffer[:m.Cursor] + "\n" + m.Buffer[m.Cursor:]
				m.Cursor++
				return m, nil
			}
			return m, interfaces.Quit
		case "backspace":
			if m.Cursor > 0 && len(m.Buffer) > 0 {
				m.Buffer = m.Buffer[:m.Cursor-1] + m.Buffer[m.Cursor:]
				m.Cursor--
			}
			return m, nil
		case "left":
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil
		case "right":
			if m.Cursor < len(m.Buffer) {
				m.Cursor++
			}
			return m, nil
		case "home":
			m.Cursor = 0
			return m, nil
		case "end":
			m.Cursor = len(m.Buffer)
			return m, nil
		default:
			if len(keyMsg.String()) == 1 && keyMsg.Type == interfaces.KeyRunes {
				m.Buffer = m.Buffer[:m.Cursor] + keyMsg.String() + m.Buffer[m.Cursor:]
				m.Cursor++
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *InputModel) UpdateWithContext(msg interfaces.Msg, ctx interfaces.Context, nav interfaces.Controller) (interfaces.Model, interfaces.Cmd) {
	return m.Update(msg)
}

func (m *InputModel) View() string {
	input := m.Buffer
	if m.Cursor >= 0 && m.Cursor <= len(input) {
		input = input[:m.Cursor] + "|" + input[m.Cursor:]
	}
	return "Input: " + strings.ReplaceAll(input, "\n", "\\n") + m.Message
}

// Observer pattern methods
func (m *InputModel) RegisterObserver(o interfaces.Observer) {
	m.observers = append(m.observers, o)
}

func (m *InputModel) UnregisterObserver(o interfaces.Observer) {
	for i, obs := range m.observers {
		if obs == o {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			break
		}
	}
}

func (m *InputModel) NotifyObservers(event interface{}) {
	for _, o := range m.observers {
		o.Notify(event)
	}
}
