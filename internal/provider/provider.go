package provider

import "context"

type Request struct {
	Expert       string
	SystemPrompt string
	Message      string
}

type Response struct {
	Text string
}

type Provider interface {
	Name() string
	Chat(context.Context, Request) (Response, error)
}
