package flows
package flows

import (
	"aichat/interfaces"
	"aichat/types"
)

// FlowItem is the interface for a single step in a flow (input, confirmation, selection, etc.)
type FlowItem interface {
	ViewState() types.ViewState
	Validate(input interface{}) error
	OnEnter(ctx types.Context, nav interfaces.Controller) error
	OnExit(ctx types.Context, nav interfaces.Controller) error
	MarshalState() ([]byte, error)
	UnmarshalState([]byte) error
}

// KeyedFlowItem extends FlowItem for steps that gather user input and propagate key:data pairs
// (e.g., input, confirmation, selection)
type KeyedFlowItem interface {
	FlowItem
	Key() string
	Value() interface{}
}

