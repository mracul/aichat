package types
package ai

type ProviderInfo struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Stream   bool   `json:"stream"`
}
