package flashcardlist

import (
	"context"
	"fmt"
	"time"

	sqlc "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
	pgxutil "github.com/beerpoldet/flash-card-go/internal/adapters/pgxutil"
)

type postgres struct {
	query *sqlc.Queries
}

func newPostgres(query *sqlc.Queries) *postgres {
	return &postgres{query: query}
}

func (p *postgres) ListFlashcardsByCreatedAtAsc(
	ctx context.Context,
	limit int32,
) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardListByCreatedAtAsc(ctx, limit)
	if err != nil {
		return []flashcard{}, fmt.Errorf("ListFlashcards: %w", err)
	}
	return asFlashcards(flashcards), nil
}

func (p *postgres) ListFlashcardsByCreatedAtAscAfter(
	ctx context.Context,
	createdAt time.Time,
	id int32,
	limit int32,
) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardListByCreatedAtAscAfter(
		ctx,
		sqlc.FlashcardListByCreatedAtAscAfterParams{
			Lim:       limit,
			CreatedAt: pgxutil.PgTimestamp(createdAt),
			ID:        id,
		},
	)
	if err != nil {
		return []flashcard{}, fmt.Errorf("ListFlashcards: %w", err)
	}
	return asFlashcards(flashcards), nil
}

func (p *postgres) ListFlashcardsByCreatedAtDesc(
	ctx context.Context,
	limit int32,
) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardListByCreatedAtDesc(ctx, limit)
	if err != nil {
		return []flashcard{}, fmt.Errorf("ListFlashcards: %w", err)
	}
	return asFlashcards(flashcards), nil
}

func (p *postgres) ListFlashcardsByCreatedAtDescAfter(
	ctx context.Context,
	createdAt time.Time,
	id int32,
	limit int32,
) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardListByCreatedAtDescAfter(
		ctx,
		sqlc.FlashcardListByCreatedAtDescAfterParams{
			Lim:       limit,
			CreatedAt: pgxutil.PgTimestamp(createdAt),
			ID:        id,
		},
	)
	if err != nil {
		return []flashcard{}, fmt.Errorf("ListFlashcards: %w", err)
	}
	return asFlashcards(flashcards), nil
}

func asFlashcards(flashcards []sqlc.Flashcard) []flashcard {
	flashcardModels := make([]flashcard, 0, len(flashcards))
	for _, f := range flashcards {
		flashcardModels = append(flashcardModels, asFlashcard(&f))
	}
	return flashcardModels
}

func asFlashcard(f *sqlc.Flashcard) flashcard {
	return flashcard{
		ID:        f.ID,
		Title:     f.Title,
		Content:   f.Content,
		CreatedAt: f.CreatedAt.Time,
	}
}
