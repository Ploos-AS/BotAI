package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ploos-AS/BotAI/internal/expert"
	"github.com/Ploos-AS/BotAI/internal/provider"
)

const maxHistory = 20
const maxMessageBytes = 4096

type Server struct{
	provider provider.Provider
	metrics metrics
}

func New(p provider.Provider) *Server { return &Server{provider: p} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /v1/experts", s.experts)
	mux.HandleFunc("GET /v1/status", s.status)
	mux.HandleFunc("GET /metrics", s.metricsHandler)
	mux.HandleFunc("POST /v1/chat", s.chat)
	return s.observe(mux)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "provider": s.provider.Name()})
}

func (s *Server) experts(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, expert.List())
}

func (s *Server) status(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"provider": s.provider.Name(),
		"experts": len(expert.List()),
	})
}

func (s *Server) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	requests := s.metrics.requests.Load()
	chats := s.metrics.chatRequests.Load()
	errors := s.metrics.chatErrors.Load()
	latency := s.metrics.chatLatencyNS.Load()
	fmt.Fprintf(w, "# TYPE botai_http_requests_total counter\nbotai_http_requests_total %d\n", requests)
	fmt.Fprintf(w, "# TYPE botai_chat_requests_total counter\nbotai_chat_requests_total %d\n", chats)
	fmt.Fprintf(w, "# TYPE botai_chat_errors_total counter\nbotai_chat_errors_total %d\n", errors)
	fmt.Fprintf(w, "# TYPE botai_chat_duration_seconds_sum counter\nbotai_chat_duration_seconds_sum %.6f\n", float64(latency)/1e9)
	fmt.Fprintf(w, "# TYPE botai_chat_duration_seconds_count counter\nbotai_chat_duration_seconds_count %d\n", chats)
}

func validText(s string) bool {
	return strings.TrimSpace(s) != "" && len(s) <= maxMessageBytes
}

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Expert  string `json:"expert"`
		History []struct {
			Role string `json:"role"`
			Content string `json:"content"`
		} `json:"history"`
		Message string `json:"message"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 128<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	in.Expert = strings.TrimSpace(in.Expert)
	in.Message = strings.TrimSpace(in.Message)
	if in.Expert == "" { in.Expert = "general" }
	if !validText(in.Message) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "message is required and must be at most 4096 bytes"})
		return
	}
	var profile expert.Profile
	var ok bool
	if in.Expert == "auto" {
		profile = expert.Route(in.Message)
	} else {
		profile, ok = expert.Get(in.Expert)
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown expert"})
			return
		}
	}
	if len(in.History) > maxHistory {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "history exceeds 20 messages"})
		return
	}
	history := make([]provider.Message, 0, len(in.History))
	for _, m := range in.History {
		role := strings.TrimSpace(m.Role)
		content := strings.TrimSpace(m.Content)
		if (role != "user" && role != "assistant") || !validText(content) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid history message"})
			return
		}
		history = append(history, provider.Message{Role:role, Content:content})
	}
	out, err := s.provider.Chat(r.Context(), provider.Request{
		Expert: profile.ID, SystemPrompt: profile.SystemPrompt, History: history, Message: in.Message,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "provider unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"expert": profile.ID, "text": out.Text, "provider": s.provider.Name()})
}
