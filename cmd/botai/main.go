package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	expertPath := strings.TrimSpace(os.Getenv("BOTAI_EXPERTS_FILE"))
	if expertPath != "" {
		if err := expert.LoadFile(expertPath); err != nil { log.Fatal(err) }
		log.Printf("loaded expert configuration from %s", expertPath)

		hup := make(chan os.Signal, 1)
		signal.Notify(hup, syscall.SIGHUP)
		go func() {
			for range hup {
				if err := expert.LoadFile(expertPath); err != nil {
					log.Printf("expert configuration reload failed; keeping current registry: %v", err)
					continue
				}
				log.Printf("reloaded expert configuration from %s", expertPath)
			}
		}()
	}

	p, err := configuredProvider()
	if err != nil { log.Fatal(err) }

	s := api.New(p)
	httpServer := &http.Server{
		Addr: addr,
		Handler: s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout: 90 * time.Second,
		MaxHeaderBytes: 16 << 10,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Printf("shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}()

	log.Printf("BotAI M0.9 listening on %s provider=%s", addr, p.Name())
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	log.Printf("BotAI stopped")
}
