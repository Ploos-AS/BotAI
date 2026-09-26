package provider

import "context"

type Message struct {
	Role    string
	Content string
}

type Request struct {
	Expert       string
	SystemPrompt string
	History      []Message
	Message      string
}

type Response struct {
	Text string
}

type Provider interface {
	Name() string
	Chat(context.Context, Request) (Response, error)
}
