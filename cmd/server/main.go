package main

import (
	"net"
	"net/http"

	"github.com/beerpoldet/flash-card-go/internal/adapters/logger"
	"github.com/beerpoldet/flash-card-go/internal/adapters/postgres"
	listflashcards "github.com/beerpoldet/flash-card-go/internal/usecase/list-flashcards"
)

func main() {
	logger := logger.NewLogger()

	mux := http.NewServeMux()

	postgresClient := postgres.NewPostgres()

	listflashcardsHandler := listflashcards.NewHandler(postgresClient)
	mux.Handle("GET /", listflashcardsHandler)

	config, err := loadConfig()
	if err != nil {
		logger.Error("loadConfig failed", "error", err)
		return
	}

	listener, err := net.Listen("tcp", config.Server.Address)
	if err != nil {
		logger.Error("listen failed", "error", err)
		return
	}
	logger.Info("Starting server at", "address", "http://"+config.Server.Address)
	err = http.Serve(listener, mux)
	if err != nil {
		logger.Error("Server stopped", "err", err)
	}
}
