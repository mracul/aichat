package ai
package ai

import "aichat/services/ai/types"

type AIProvider interface {
	Info() types.ProviderInfo
	SendMessage(messages []map[string]string, apiKey, model string) (string, error)
	StreamMessage(messages []map[string]string, apiKey, model string, onData func(data string)) error
}

