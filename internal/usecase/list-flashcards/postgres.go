package listflashcards

import (
	"context"

	sqlc "github.com/beerpoldet/flash-card-go/gen/sqlc/flashcardlist"
)

type dbClient struct {
	query *sqlc.Queries
}

func NewDBClient(query *sqlc.Queries) *dbClient {
	return &dbClient{query: query}
}

func (p *dbClient) ListFlashcards(ctx context.Context, limit int32) ([]flashcard, error) {
	flashcards, err := p.query.FlashcardList(ctx, limit)
	flashcardModels := []flashcard{}
	if err != nil {
		return flashcardModels, err
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
