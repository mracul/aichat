package flows
package flows

import (
	"fmt"
)

// FlowStruct orchestrates a sequence of FlowItems, collects their data, and triggers OnSuccess when complete.
type FlowStruct struct {
	RequiredKeys []string
	Items        []FlowItem
	Data         map[string]string
	OnSuccess    func(data map[string]string) error
}

// SetData sets a key-value pair in the flow's data map.
func (f *FlowStruct) SetData(key, value string) {
	if f.Data == nil {
		f.Data = make(map[string]string)
	}
	f.Data[key] = value
}

// MergeChildData merges a child flow's data into this flow's data.
func (f *FlowStruct) MergeChildData(child map[string]string) {
	for k, v := range child {
		f.SetData(k, v)
	}
}

// HasKey checks if a key is present in the data map.
func (f *FlowStruct) HasKey(k string) bool {
	_, ok := f.Data[k]
	return ok
}

// MissingKeys returns a slice of required keys that are not yet present.
func (f *FlowStruct) MissingKeys() []string {
	var missing []string
	for _, k := range f.RequiredKeys {
		if !f.HasKey(k) {
			missing = append(missing, k)
		}
	}
	return missing
}

// Run executes each FlowItem, collects data, and calls OnSuccess if complete.
func (f *FlowStruct) Run() error {
	for _, item := range f.Items {
		// item.Submit(f) is invalid; FlowItem does not have Submit
		// If needed, call OnEnter or OnExit, or process item as appropriate
		if err := item.OnExit(nil, nil); err != nil { // TODO: pass correct ctx, nav if needed
			return err
		}
	}
	if len(f.MissingKeys()) == 0 {
		if f.OnSuccess != nil {
			return f.OnSuccess(f.Data)
		}
		return nil
	}
	return fmt.Errorf("incomplete data: missing %v", f.MissingKeys())
}

