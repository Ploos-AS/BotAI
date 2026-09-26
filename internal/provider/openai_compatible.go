package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAICompatible struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

func NewOpenAICompatible(baseURL, apiKey, model string) (*OpenAICompatible, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	model = strings.TrimSpace(model)
	if baseURL == "" { return nil, errors.New("provider base URL is required") }
	if model == "" { return nil, errors.New("provider model is required") }
	return &OpenAICompatible{
		baseURL: baseURL, apiKey: strings.TrimSpace(apiKey), model: model,
		client: &http.Client{Timeout: 45 * time.Second},
	}, nil
}

func (p *OpenAICompatible) Name() string { return "openai-compatible" }

func (p *OpenAICompatible) Chat(ctx context.Context, r Request) (Response, error) {
	body := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are the " + r.Expert + " expert for an IRC bot. Be concise, accurate, and IRC-friendly."},
			{"role": "user", "content": r.Message},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil { return Response{}, err }
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil { return Response{}, err }
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" { req.Header.Set("Authorization", "Bearer "+p.apiKey) }

	resp, err := p.client.Do(req)
	if err != nil { return Response{}, err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
		return Response{}, fmt.Errorf("provider returned HTTP %d", resp.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct { Content string `json:"content"` } `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&out); err != nil { return Response{}, err }
	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		return Response{}, errors.New("provider returned no message")
	}
	return Response{Text: strings.TrimSpace(out.Choices[0].Message.Content)}, nil
}
