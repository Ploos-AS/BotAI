package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ploos-AS/BotAI/internal/api"
	"github.com/Ploos-AS/BotAI/internal/expert"
	"github.com/Ploos-AS/BotAI/internal/provider"
)

func configuredProvider() (provider.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("BOTAI_PROVIDER"))) {
	case "", "echo":
		return provider.NewEcho(), nil
	case "openai-compatible":
		return provider.NewOpenAICompatible(
			os.Getenv("BOTAI_BASE_URL"),
			os.Getenv("BOTAI_API_KEY"),
			os.Getenv("BOTAI_MODEL"),
		)
	default:
		return nil, fmt.Errorf("unknown BOTAI_PROVIDER")
	}
}

func main() {
	addr := os.Getenv("BOTAI_LISTEN")
	if addr == "" { addr = "127.0.0.1:8090" }

	if path := strings.TrimSpace(os.Getenv("BOTAI_EXPERTS_FILE")); path != "" {
		if err := expert.LoadFile(path); err != nil { log.Fatal(err) }
		log.Printf("loaded expert configuration from %s", path)
	}

	p, err := configuredProvider()
	if err != nil { log.Fatal(err) }

	s := api.New(p)
	log.Printf("BotAI M0.6 listening on %s provider=%s", addr, p.Name())
	if err := http.ListenAndServe(addr, s.Handler()); err != nil { log.Fatal(err) }
}
