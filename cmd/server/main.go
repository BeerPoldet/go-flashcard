package main

import (
	"context"
	"net"
	"net/http"

	"github.com/beerpoldet/flash-card-go/internal/adapters/logger"
	flashcardlist "github.com/beerpoldet/flash-card-go/internal/usecase/flashcardlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := logger.NewLogger()

	config, err := loadConfig()
	if err != nil {
		logger.Error("loadConfig failed", "error", err)
		return
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.DB.ConnectionURI)
	if err != nil {
		logger.Error("create pgxpool failed", "error", err)
		return
	}
	defer pool.Close()

	mux := http.NewServeMux()
	flashcardlistHandler := flashcardlist.NewLiveHandler(pool, logger)
	mux.Handle("GET /flashcards", flashcardlistHandler)

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
