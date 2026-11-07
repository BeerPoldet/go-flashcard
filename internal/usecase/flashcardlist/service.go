package flashcardlist

import "context"

type db interface {
	ListFlashcards(ctx context.Context, limit int32) ([]flashcard, error)
}

type service struct {
	db     db
	logger logger
}

func newService(db db, logger logger) *service {
	return &service{
		db:     db,
		logger: logger,
	}
}

func (uc *service) Perform(ctx context.Context, limit int) ([]flashcard, error) {
	return uc.db.ListFlashcards(ctx, int32(limit))
}
