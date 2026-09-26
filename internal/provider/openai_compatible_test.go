package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAICompatibleChat(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" { t.Fatalf("path=%s", r.URL.Path) }
		if got := r.Header.Get("Authorization"); got != "Bearer secret" { t.Fatalf("auth=%q", got) }
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"433 means nickname in use."}}]}`))
	}))
	defer s.Close()
	p, err := NewOpenAICompatible(s.URL+"/v1", "secret", "test-model")
	if err != nil { t.Fatal(err) }
	got, err := p.Chat(context.Background(), Request{Expert:"irc", Message:"What is 433?"})
	if err != nil { t.Fatal(err) }
	if got.Text != "433 means nickname in use." { t.Fatalf("text=%q", got.Text) }
}

func TestOpenAICompatibleRequiresConfig(t *testing.T) {
	if _, err := NewOpenAICompatible("", "", "model"); err == nil { t.Fatal("expected base URL error") }
	if _, err := NewOpenAICompatible("http://localhost", "", ""); err == nil { t.Fatal("expected model error") }
}
