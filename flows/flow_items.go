package flows
package flows

import (
	"aichat/interfaces"
	"aichat/types"
)

// InputFlowItem prompts the user for input and stores the result under a key
// Implements KeyedFlowItem

type InputFlowItem struct {
	key    string
	prompt string
	value  string
}

func NewInputFlowItem(key, prompt string) *InputFlowItem {
	return &InputFlowItem{key: key, prompt: prompt}
}

func (i *InputFlowItem) Key() string                                                { return i.key }
func (i *InputFlowItem) Value() interface{}                                         { return i.value }
func (i *InputFlowItem) ViewState() types.ViewState                                 { return nil } // TODO: implement
func (i *InputFlowItem) Validate(input interface{}) error                           { return nil } // TODO: implement
func (i *InputFlowItem) OnEnter(ctx types.Context, nav interfaces.Controller) error { return nil }
func (i *InputFlowItem) OnExit(ctx types.Context, nav interfaces.Controller) error  { return nil }
func (i *InputFlowItem) MarshalState() ([]byte, error)                              { return nil, nil }
func (i *InputFlowItem) UnmarshalState([]byte) error                                { return nil }

// ConfirmationFlowItem prompts the user for confirmation (yes/no) and stores the result under a key
// Implements KeyedFlowItem

type ConfirmationFlowItem struct {
	key     string
	message string
	value   bool
}

func NewConfirmationFlowItem(key, message string) *ConfirmationFlowItem {
	return &ConfirmationFlowItem{key: key, message: message}
}

func (c *ConfirmationFlowItem) Key() string                      { return c.key }
func (c *ConfirmationFlowItem) Value() interface{}               { return c.value }
func (c *ConfirmationFlowItem) ViewState() types.ViewState       { return nil } // TODO: implement
func (c *ConfirmationFlowItem) Validate(input interface{}) error { return nil } // TODO: implement
func (c *ConfirmationFlowItem) OnEnter(ctx types.Context, nav interfaces.Controller) error {
	return nil
}
func (c *ConfirmationFlowItem) OnExit(ctx types.Context, nav interfaces.Controller) error { return nil }
func (c *ConfirmationFlowItem) MarshalState() ([]byte, error)                             { return nil, nil }
func (c *ConfirmationFlowItem) UnmarshalState([]byte) error                               { return nil }

// NoticeFlowItem displays a notice or message to the user
type NoticeFlowItem struct {
	Message string
}

func (n *NoticeFlowItem) ViewState() types.ViewState                                 { return nil } // TODO: implement
func (n *NoticeFlowItem) Validate(input interface{}) error                           { return nil }
func (n *NoticeFlowItem) OnEnter(ctx types.Context, nav interfaces.Controller) error { return nil }
func (n *NoticeFlowItem) OnExit(ctx types.Context, nav interfaces.Controller) error  { return nil }
func (n *NoticeFlowItem) MarshalState() ([]byte, error)                              { return nil, nil }
func (n *NoticeFlowItem) UnmarshalState([]byte) error                                { return nil }

