package ai
package ai

import (
	"aichat/services/ai/providers"
	aitypes "aichat/services/ai/types"
	"encoding/json"
	"io/ioutil"
)

var providerRegistry = map[string]AIProvider{}

func LoadProvidersFromJSON(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var entries []aitypes.ProviderInfo
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for _, entry := range entries {
		var p AIProvider
		if entry.Name == "OpenAI" || entry.Name == "OpenAI (s)" {
			p = providers.NewOpenAIProvider(entry.Stream)
		} else if entry.Name == "OpenRouter" || entry.Name == "OpenRouter (s)" {
			p = providers.NewOpenRouterProvider(entry.Stream)
		}
		if p != nil {
			providerRegistry[entry.Name] = p
		}
	}
	return nil
}

func GetAllProviders() []AIProvider {
	providers := []AIProvider{}
	for _, p := range providerRegistry {
		providers = append(providers, p)
	}
	return providers
}

func GetProviderByName(name string) AIProvider {
	return providerRegistry[name]
}
