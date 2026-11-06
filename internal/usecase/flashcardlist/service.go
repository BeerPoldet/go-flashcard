package flashcardlist

import "context"

type useCase struct {
	db     db
	logger logger
	output useCaseOutputPort
}

func NewUseCase(db db, logger logger, output useCaseOutputPort) *useCase {
	return &useCase{
		db:     db,
		logger: logger,
		output: output,
	}
}

func (uc *useCase) listFlashcards(ctx context.Context, limit int) {
	uc.output.getFlashcardListResponse(uc.db.ListFlashcards(ctx, int32(limit)))
}

type useCaseInputPort interface {
	listFlashcards(context.Context, int)
}

type useCaseOutputPort interface {
	getFlashcardListResponse([]flashcard, error)
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
