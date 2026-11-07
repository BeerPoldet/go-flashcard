package flashcardlist

import (
	"context"
	"fmt"

	sqlc "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
)

type postgres struct {
	query *sqlc.Queries
}

func newPostgres(query *sqlc.Queries) *postgres {
	return &postgres{query: query}
}

func (p *postgres) ListFlashcards(ctx context.Context, limit int32) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardList(ctx, limit)
	flashcardModels := []flashcard{}
	if err != nil {
		return flashcardModels, fmt.Errorf("ListFlashcards: %w", err)
	}
	for _, f := range flashcards {
		flashcardModels = append(flashcardModels, asFlashcard(&f))
	}
	return flashcardModels, nil
}

func asFlashcard(f *sqlc.Flashcard) flashcard {
	return flashcard{
		ID:        f.ID,
		Title:     f.Title,
		Content:   f.Content,
		CreatedAt: f.CreatedAt.Time,
	}
}
