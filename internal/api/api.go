package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ploos-AS/BotAI/internal/expert"
	"github.com/Ploos-AS/BotAI/internal/provider"
)

type Server struct{ provider provider.Provider }

func New(p provider.Provider) *Server { return &Server{provider: p} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /v1/experts", s.experts)
	mux.HandleFunc("POST /v1/chat", s.chat)
	return mux
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

func (s *Server) chat(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Expert  string `json:"expert"`
		Message string `json:"message"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	in.Expert = strings.TrimSpace(in.Expert)
	in.Message = strings.TrimSpace(in.Message)
	if in.Expert == "" { in.Expert = "general" }
	if _, ok := expert.Get(in.Expert); !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown expert"})
		return
	}
	if in.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "message is required"})
		return
	}
	out, err := s.provider.Chat(r.Context(), provider.Request{Expert: in.Expert, Message: in.Message})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "provider unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"expert": in.Expert, "text": out.Text, "provider": s.provider.Name()})
}
