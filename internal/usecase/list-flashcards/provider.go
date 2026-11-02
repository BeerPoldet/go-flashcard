package listflashcards

import (
	flashcardlistdb "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLiveHandler(pool *pgxpool.Pool, logger logger) *Handler {
	queries := flashcardlistdb.New(pool)
	db := NewDBClient(queries)
	svc := NewService(db, logger)
	return &Handler{svc: svc}
}
