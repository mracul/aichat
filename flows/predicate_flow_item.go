package flows
package flows

import (
	"aichat/interfaces"
	"aichat/types"
)

// PredicateFlowItem is the single branching flow item for all conditional flows.
// It executes YesFlow if Predicate returns true, otherwise NoFlow.
type PredicateFlowItem struct {
	Predicate func(ctx types.Context) bool
	YesFlow   FlowItem // Can be a *Flow or any FlowItem
	NoFlow    FlowItem // Can be a *Flow or any FlowItem
}

func (f *PredicateFlowItem) Run(ctx types.Context, nav interfaces.Controller) (FlowResult, error) {
	if f.Predicate(ctx) {
		if f.YesFlow != nil {
			if err := f.YesFlow.OnEnter(ctx, nav); err != nil {
				return FlowResult{}, err
			}
			return FlowResult{Value: f.YesFlow, Success: true}, nil
		}
	} else {
		if f.NoFlow != nil {
			if err := f.NoFlow.OnEnter(ctx, nav); err != nil {
				return FlowResult{}, err
			}
			return FlowResult{Value: f.NoFlow, Success: true}, nil
		}
	}
	return FlowResult{}, nil
}

// Example usage:
// yesFlow := NewFlow([]FlowItem{/* ... */}, handlerForYes)
// noFlow := NewFlow([]FlowItem{/* ... */}, handlerForNo)
// fork := &PredicateFlowItem{
//     Predicate: func(ctx types.Context) bool { return ctx["shouldProceed"].(bool) },
//     YesFlow:   yesFlow,
//     NoFlow:    noFlow,
// }
//
// This fork does not show any notice, just branches to the correct flow.

