package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hetzner/cert-manager-webhook-hetzner/internal/hetzner"
	"github.com/hetzner/cert-manager-webhook-hetzner/internal/version"
)

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	logger := newLogger()

	logger.Info("starting webhook", "version", version.Version)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewBuildInfoCollector())

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("starting metrics server", "addr", server.Addr)
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	// This will register our custom DNS provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.
	cmd.RunWebhookServer(GroupName, hetzner.New(logger, registry))

	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		logger.Info("shutting down metrics server")

		if err := server.Shutdown(ctx); err != nil {
			logger.Error(err.Error())
		}

		cancel()
	}
}
