package flows
package flows

import (
	"fmt"
)

// Global/common flows

// ExitFlow prompts the user for confirmation before exiting.
func ExitFlow(onExit func() error) *FlowStruct {
	return &FlowStruct{
		RequiredKeys: []string{"confirm_exit"},
		Data:         make(map[string]string),
		Items:        []FlowItem{},
		OnSuccess: func(data map[string]string) error {
			if v, ok := data["confirm_exit"]; ok && (v == "yes" || v == "y") {
				if onExit != nil {
					return onExit()
				}
				return nil
			}
			return fmt.Errorf("exit cancelled")
		},
	}
}

