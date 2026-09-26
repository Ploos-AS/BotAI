package provider

import "context"

// Echo is the deterministic M0 provider. It proves the API/provider boundary
// without requiring a cloud account, model download, or external service.
type Echo struct{}

func NewEcho() Echo { return Echo{} }
func (Echo) Name() string { return "echo" }
func (Echo) Chat(_ context.Context, r Request) (Response, error) {
	return Response{Text: "[" + r.Expert + "] " + r.Message}, nil
}
