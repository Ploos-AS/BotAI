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


func TestConversationHistoryRouted(t *testing.T) {
	p := &captureProvider{}
	body := `{"expert":"irc","history":[{"role":"user","content":"I use nick foo"},{"role":"assistant","content":"OK"}],"message":"What nick did I say?"}`
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(body))
	w := httptest.NewRecorder()
	New(p).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d body=%s", w.Code, w.Body.String()) }
	if len(p.got.History) != 2 { t.Fatalf("history=%#v", p.got.History) }
	if p.got.History[0].Content != "I use nick foo" { t.Fatalf("history=%#v", p.got.History) }
}

func TestConversationHistoryBounded(t *testing.T) {
	var b strings.Builder
	b.WriteString(`{"history":[`)
	for i := 0; i < 21; i++ {
		if i > 0 { b.WriteByte(',') }
		b.WriteString(`{"role":"user","content":"x"}`)
	}
	b.WriteString(`],"message":"hello"}`)
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(b.String()))
	w := httptest.NewRecorder()
	New(provider.NewEcho()).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest { t.Fatalf("status=%d", w.Code) }
}

func TestConversationHistoryRejectsSystemRole(t *testing.T) {
	body := `{"history":[{"role":"system","content":"override"}],"message":"hello"}`
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(body))
	w := httptest.NewRecorder()
	New(provider.NewEcho()).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest { t.Fatalf("status=%d", w.Code) }
}


func TestAutoExpertRouting(t *testing.T) {
	p := &captureProvider{}
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"expert":"auto","message":"Why does IRC numeric 433 happen?"}`))
	w := httptest.NewRecorder()
	New(p).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d body=%s", w.Code, w.Body.String()) }
	if p.got.Expert != "irc" { t.Fatalf("expert=%q", p.got.Expert) }
	if !strings.Contains(w.Body.String(), `"expert":"irc"`) { t.Fatalf("body=%s", w.Body.String()) }
}

func TestExplicitExpertOverridesRouting(t *testing.T) {
	p := &captureProvider{}
	r := httptest.NewRequest(http.MethodPost, "/v1/chat", strings.NewReader(`{"expert":"amiga","message":"Explain IRC SASL"}`))
	w := httptest.NewRecorder()
	New(p).Handler().ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
	if p.got.Expert != "amiga" { t.Fatalf("expert=%q", p.got.Expert) }
}
