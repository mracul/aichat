package interfaces
package interfaces

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Controller interface {
	Push(view interface{})
	Pop() interface{}
	Replace(view interface{})
	Current() interface{}
	CanPop() bool
	ShowModal(modalType string, data interface{})
	HideModal()
}

type ViewType int

const (
	MenuStateType ViewType = iota
	ChatStateType
	ModalStateType
)

type ControlType struct {
	Name   string
	Key    tea.KeyType
	Action func() bool
}

type ControlSet struct {
	Controls []ControlType
}

type ControlSetProvider interface {
	GetControlSet(ctx Context) ControlSet
}

type Observer interface {
	Notify(event interface{})
}

type Renderable interface {
	View() string
}

type Updatable interface {
	Update(msg Msg) (Model, Cmd)
	UpdateWithContext(msg Msg, ctx Context, nav Controller) (Model, Cmd)
	Init() Cmd
}

type Cmd = tea.Cmd

type Msg = tea.Msg

type Model = tea.Model

// Key aliases
var (
	KeyEnter = tea.KeyEnter
	KeyCtrlV = tea.KeyCtrlV
	KeyCtrlX = tea.KeyCtrlX
	KeyCtrlC = tea.KeyCtrlC
	KeyEsc   = tea.KeyEsc
	KeyCtrlI = tea.KeyCtrlI
	KeyRunes = tea.KeyRunes
	Quit     = tea.Quit
)

type KeyMsg = tea.KeyMsg

