package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ploos-AS/BotAI/internal/provider"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	New(provider.NewEcho()).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
}

func TestChatIRCExpert(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"expert":"irc","message":"hello"}`))
	w := httptest.NewRecorder()
	New(provider.NewEcho()).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d body=%s", w.Code, w.Body.String()) }
	if !strings.Contains(w.Body.String(), "[irc] hello") { t.Fatalf("body=%s", w.Body.String()) }
}

func TestUnknownExpert(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"expert":"nope","message":"hello"}`))
	w := httptest.NewRecorder()
	New(provider.NewEcho()).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest { t.Fatalf("status=%d", w.Code) }
}

type captureProvider struct{ got provider.Request }
func (p *captureProvider) Name() string { return "capture" }
func (p *captureProvider) Chat(_ context.Context, r provider.Request) (provider.Response, error) {
	p.got = r
	return provider.Response{Text:"ok"}, nil
}

func TestExpertPromptResolvedBeforeProvider(t *testing.T) {
	p := &captureProvider{}
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"expert":"amiga","message":"ARexx?"}`))
	w := httptest.NewRecorder()
	New(p).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
	if p.got.Expert != "amiga" { t.Fatalf("expert=%q", p.got.Expert) }
	if !strings.Contains(p.got.SystemPrompt, "Amiga") || !strings.Contains(p.got.SystemPrompt, "ARexx") {
		t.Fatalf("unexpected prompt=%q", p.got.SystemPrompt)
	}
}
