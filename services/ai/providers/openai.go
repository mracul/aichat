package providers
package providers

import (
	"aichat/services/ai/types"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type OpenAIProvider struct {
	info types.ProviderInfo
}

func NewOpenAIProvider(stream bool) *OpenAIProvider {
	endpoint := "https://api.openai.com/v1/chat/completions"
	name := "OpenAI"
	if stream {
		name += " (s)"
	}
	return &OpenAIProvider{
		info: types.ProviderInfo{
			Name:     name,
			Endpoint: endpoint,
			Stream:   stream,
		},
	}
}

func (p *OpenAIProvider) Info() types.ProviderInfo {
	return p.info
}

func (p *OpenAIProvider) SendMessage(messages []map[string]string, apiKey, model string) (string, error) {
	body := map[string]interface{}{
		"model":    model,
		"messages": messages,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", p.info.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no choices in response")
	}
	msg, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", errors.New("no message content in response")
	}
	return msg, nil
}

func (p *OpenAIProvider) StreamMessage(messages []map[string]string, apiKey, model string, onData func(data string)) error {
	body := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   true,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", p.info.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := dec.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if data, ok := chunk["choices"].([]interface{}); ok && len(data) > 0 {
			delta := data[0].(map[string]interface{})["delta"].(map[string]interface{})
			if content, ok := delta["content"].(string); ok {
				onData(content)
			}
		}
	}
	return nil
}

