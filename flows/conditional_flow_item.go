package flows
package flows

// ConditionalFlowItem branches to YesPath or NoPath based on a condition, and merges the chosen subflow's data into the parent.
type ConditionalFlowItem struct {
	key       string
	condition func() bool
	YesPath   *FlowStruct
	NoPath    *FlowStruct
}

func NewConditionalFlowItem(key string, condition func() bool, yesPath, noPath *FlowStruct) *ConditionalFlowItem {
	return &ConditionalFlowItem{
		key:       key,
		condition: condition,
		YesPath:   yesPath,
		NoPath:    noPath,
	}
}

func (c *ConditionalFlowItem) Key() string { return c.key }

// Gather returns the string representation of the branch taken ("yes" or "no").
func (c *ConditionalFlowItem) Gather() (string, error) {
	if c.condition() {
		return "yes", nil
	}
	return "no", nil
}

// Submit runs the chosen subflow and merges its data into the parent.
func (c *ConditionalFlowItem) Submit(parent *FlowStruct) error {
	var chosen *FlowStruct
	var branch string
	if c.condition() {
		chosen = c.YesPath
		branch = "yes"
	} else {
		chosen = c.NoPath
		branch = "no"
	}
	// Record which branch was taken
	parent.SetData(c.key, branch)
	if chosen != nil {
		err := chosen.Run()
		if err != nil {
			return err
		}
		parent.MergeChildData(chosen.Data)
	}
	return nil
}

