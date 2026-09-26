package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ploos-AS/BotAI/internal/api"
	"github.com/Ploos-AS/BotAI/internal/provider"
)

func main() {
	addr := os.Getenv("BOTAI_LISTEN")
	if addr == "" {
		addr = "127.0.0.1:8090"
	}

	s := api.New(provider.NewEcho())
	log.Printf("BotAI M0 listening on %s", addr)
	if err := http.ListenAndServe(addr, s.Handler()); err != nil {
		log.Fatal(err)
	}
}
