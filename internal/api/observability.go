package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type metrics struct {
	requests atomic.Uint64
	chatRequests atomic.Uint64
	chatErrors atomic.Uint64
	chatLatencyNS atomic.Uint64
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" && len(id) <= 128 { return id }
	var b [8]byte
	if _, err := rand.Read(b[:]); err == nil { return hex.EncodeToString(b[:]) }
	return "unknown"
}

func (s *Server) observe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := requestID(r)
		w.Header().Set("X-Request-ID", id)
		sw := &statusWriter{ResponseWriter:w, status:http.StatusOK}
		s.metrics.requests.Add(1)
		if r.Method == http.MethodPost && r.URL.Path == "/v1/chat" { s.metrics.chatRequests.Add(1) }
		next.ServeHTTP(sw, r)
		d := time.Since(start)
		if r.Method == http.MethodPost && r.URL.Path == "/v1/chat" {
			s.metrics.chatLatencyNS.Add(uint64(d))
			if sw.status >= 400 { s.metrics.chatErrors.Add(1) }
		}
		log.Printf("request id=%s method=%s path=%s status=%d duration_ms=%d", id, r.Method, r.URL.Path, sw.status, d.Milliseconds())
	})
}
