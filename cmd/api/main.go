package main

import (
	"log"
	"net/http"

	adapterhttp "github.com/willejs/ports-service/internal/adapter/http"
	"github.com/willejs/ports-service/internal/adapter/repository"
	"github.com/willejs/ports-service/internal/app"
	"github.com/willejs/ports-service/internal/controller"
	"github.com/willejs/ports-service/internal/infrastructure/memdb"
)

func main() {
	// Initialize MemDB
	db, err := memdb.NewMemDB()
	if err != nil {
		log.Fatalf("error initializing MemDB: %v", err)
	}

	// Create repository, service, controller, and handler
	portRepo := repository.NewMemDBPortRepository(db)
	portService := app.NewPortService(portRepo)
	portController := controller.NewPortController(portService)
	portHandler := adapterhttp.NewPortHandler(portController)

	// Load ports data from JSON file
	log.Println("Upserting ports from file...")
	if err := portController.UpsertPortsFromFile(); err != nil {
		log.Fatalf("error upserting ports from file: %v", err)
	}

	mux := http.NewServeMux()

	// Set up HTTP server and routes
	mux.HandleFunc("/ports", portHandler.ListPorts)
	mux.HandleFunc("/ready", portHandler.Ready)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting server: %v", err)
	}
}
