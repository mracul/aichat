package flows
package flows

import (
	"aichat/interfaces"
	"aichat/types"
)

// FlowResult is a generic result from a FlowItem
// Can be extended with more fields as needed
type FlowResult struct {
	Value   interface{} // The value/result of the flow item
	Success bool        // Indicates if the flow item succeeded
	Err     error       // Error, if any
}

// Flow is the interface for orchestrating multi-step processes (flows)
type Flow interface {
	Start(ctx types.Context, nav interfaces.Controller) error
	Next(msg interface{}, ctx types.Context, nav interfaces.Controller) error
	Cancel(ctx types.Context, nav interfaces.Controller) error
	Pause() ([]byte, error)
	Resume(data []byte, ctx types.Context, nav interfaces.Controller) error
	CurrentItem() FlowItem
}

// FlowImpl manages a sequence of FlowItems, supporting conditional branching
// Now supports context data gathering and a final handler function
type FlowImpl struct {
	Items   []FlowItem
	Context map[string]interface{}                                            // Stores key:value data from FlowItems
	Handler func(ctx map[string]interface{}, nav interfaces.Controller) error // Called after all items
	current int
}

func NewFlowImpl(items []FlowItem, handler func(ctx map[string]interface{}, nav interfaces.Controller) error) *FlowImpl {
	return &FlowImpl{
		Items:   items,
		Context: make(map[string]interface{}),
		Handler: handler,
		current: 0,
	}
}

func (f *FlowImpl) Start(ctx types.Context, nav interfaces.Controller) error {
	f.current = 0
	if len(f.Items) == 0 {
		return nil
	}
	// Optionally push the first item's view state here
	return nil
}

func (f *FlowImpl) Next(msg interface{}, ctx types.Context, nav interfaces.Controller) error {
	if f.current >= len(f.Items) {
		return nil
	}
	item := f.Items[f.current]
	// If the item is a KeyedFlowItem, collect its key:value
	if keyed, ok := item.(KeyedFlowItem); ok {
		f.Context[keyed.Key()] = keyed.Value()
	}
	// Optionally process msg and advance
	f.current++
	return nil
}

func (f *FlowImpl) Cancel(ctx types.Context, nav interfaces.Controller) error {
	// Optionally pop view states or clean up
	return nil
}

func (f *FlowImpl) Pause() ([]byte, error) {
	// Serialize current state
	return nil, nil
}

func (f *FlowImpl) Resume(data []byte, ctx types.Context, nav interfaces.Controller) error {
	// Deserialize state
	return nil
}

func (f *FlowImpl) CurrentItem() FlowItem {
	if f.current < len(f.Items) {
		return f.Items[f.current]
	}
	return nil
}

