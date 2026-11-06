package flashcardlist

import (
	flashcardlistdb "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLiveHandler(pool *pgxpool.Pool, logger logger) *Handler {
	queries := flashcardlistdb.New(pool)
	db := NewDBClient(queries)
	output := newPresenter(logger)
	uc := NewUseCase(db, logger, output)
	return &Handler{uc: uc}
}
