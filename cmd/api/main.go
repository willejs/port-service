package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	adapterhttp "github.com/willejs/ports-service/internal/adapter/http"
	"github.com/willejs/ports-service/internal/adapter/repository"
	"github.com/willejs/ports-service/internal/app"
	"github.com/willejs/ports-service/internal/controller"
	infrahttp "github.com/willejs/ports-service/internal/infrastructure/http"
	"github.com/willejs/ports-service/internal/infrastructure/memdb"
	"github.com/willejs/ports-service/internal/infrastructure/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize OpenTelemetry (tracing and metrics)
	_, otelCleanup, err := otel.NewProviders("port-service")
	if err != nil {
		logger.Error("Failed to setup otel providers", slog.String("component", "main"), slog.Any("error", err))
	}
	defer otelCleanup()

	// Initialize MemDB
	db, err := memdb.NewMemDB(logger)
	if err != nil {
		logger.Error("Failed to initalize memdb", slog.String("component", "main"), slog.Any("error", err))
		os.Exit(1)
	}

	// Create repository, service, controller, and handler
	portRepo := repository.NewMemDBPortRepository(db)
	portService := app.NewPortService(portRepo)
	portController := controller.NewPortController(logger, portService)
	portHandler := adapterhttp.NewPortHandler(portController)

	// Load ports data from JSON file
	if err := portController.UpsertPortsFromFile(); err != nil {
		logger.Error("Failed to upsert all ports from file", slog.String("component", "main"), slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Set up HTTP server and routes
	mux.HandleFunc("/ports", portHandler.ListPorts)
	mux.HandleFunc("/ready", portHandler.Ready)

	// add the prom http handler so that we can scrape metrics
	mux.Handle("/metrics", promhttp.Handler())

	loggedMux := infrahttp.LoggingMiddleware(logger)(mux)
	otelLoggedMux := otelhttp.NewHandler(loggedMux, "server")

	server := &http.Server{
		Addr:    ":8080",
		Handler: otelLoggedMux,
	}

	go func() {
		logger.Info("Starting server", slog.String("component", "main"), slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	// Shut down the server gracefully
	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-stop
	logger.Info("Shutting down server", slog.String("component", "main"))

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown server", slog.String("component", "main"), slog.Any("error", err))
		os.Exit(1)
	}

	cancel()

	logger.Info("Server shutdown complete", slog.String("component", "main"))
}
