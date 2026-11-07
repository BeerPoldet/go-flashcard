package flashcardlist

import (
	flashcardlistdb "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewLiveHandler(pool *pgxpool.Pool, logger logger) *Handler {
	queries := flashcardlistdb.New(pool)
	db := newPostgres(queries)
	uc := newService(db, logger)
	return newHandler(uc, logger)
}
