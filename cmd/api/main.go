package main

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	adapterhttp "github.com/willejs/ports-service/internal/adapter/http"
	"github.com/willejs/ports-service/internal/adapter/repository"
	"github.com/willejs/ports-service/internal/app"
	"github.com/willejs/ports-service/internal/controller"
	"github.com/willejs/ports-service/internal/infrastructure/memdb"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize MemDB
	db, err := memdb.NewMemDB()
	if err != nil {
		logger.Error("Failed to initialize MemDB", slog.String("component", "main"), slog.Any("error", err))
		os.Exit(1)
	}

	// Create repository, service, controller, and handler
	portRepo := repository.NewMemDBPortRepository(db)
	portService := app.NewPortService(portRepo)
	portController := controller.NewPortController(portService)
	portHandler := adapterhttp.NewPortHandler(portController)

	// Load ports data from JSON file
	if err := portController.UpsertPortsFromFile(); err != nil {
		logger.Error("Failed to upsert all ports from file", slog.String("component", "main"), slog.Any("error", err))
		// theres no fatal level for slog
		// I would probably use a different logger in future, or just wrap the logger and create a fatal level
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Set up HTTP server and routes
	mux.HandleFunc("/ports", portHandler.ListPorts)
	mux.HandleFunc("/ready", portHandler.Ready)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	logger.Info("Starting server", slog.String("component", "main"), slog.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", slog.String("component", "main"), slog.Any("error", err))
		os.Exit(1)
	}
}
