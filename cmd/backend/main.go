package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/api"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: api.NewServer(cfg).Handler(),
	}

	log.Printf(
		"demo backend listening on %s service=%s instance=%s",
		server.Addr,
		cfg.Service,
		cfg.Instance,
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}
