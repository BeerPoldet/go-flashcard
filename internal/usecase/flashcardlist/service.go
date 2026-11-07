package flashcardlist

import "context"

type service struct {
	db     db
	logger logger
}

func NewService(db db, logger logger) *service {
	return &service{
		db:     db,
		logger: logger,
	}
}

func (uc *service) Perform(ctx context.Context, limit int) ([]flashcard, error) {
	return uc.db.ListFlashcards(ctx, int32(limit))
}

type useCase interface {
	Perform(context.Context, int) ([]flashcard, error)
}

type db interface {
	ListFlashcards(ctx context.Context, limit int32) ([]flashcard, error)
}

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
