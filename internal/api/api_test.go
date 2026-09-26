package api

import (
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
