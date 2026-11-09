package flashcardlist

import (
	"context"
	"fmt"
	"time"
)

type db interface {
	ListFlashcardsByCreatedAtAsc(ctx context.Context, limit int32) ([]flashcard, error)
	ListFlashcardsByCreatedAtAscAfter(
		ctx context.Context,
		createdAt time.Time,
		id int32,
		limit int32,
	) ([]flashcard, error)
	ListFlashcardsByCreatedAtDesc(ctx context.Context, limit int32) ([]flashcard, error)
	ListFlashcardsByCreatedAtDescAfter(
		ctx context.Context,
		createdAt time.Time,
		id int32,
		limit int32,
	) ([]flashcard, error)
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

type queryOptions struct {
	limit  int32
	cursor cursor
}

func (uc *service) Perform(
	ctx context.Context,
	options queryOptions,
) (listFlashcardQueryResult, error) {
	if options.cursor.Sort == "" {
		options.cursor.Sort = createdAtSort
	}
	if options.cursor.Order == "" {
		options.cursor.Order = descOrder
	}
	flashcards, err := uc.listFlashcards(ctx, options.cursor, options.limit+1)
	if err != nil {
		return listFlashcardQueryResult{}, err
	}
	if int32(len(flashcards)) <= options.limit { // no need to trim out the +1 limit fetched
		return listFlashcardQueryResult{data: flashcards}, nil
	}

	data := flashcards[:options.limit]
	after, err := afterCursorFrom(
		options.cursor.Sort,
		options.cursor.Order,
		data[len(data)-1],
	)
	if err != nil {
		return listFlashcardQueryResult{}, err
	}
	return listFlashcardQueryResult{
		data:  data,
		after: after,
	}, nil
}

func (uc *service) listFlashcards(
	ctx context.Context,
	cur cursor,
	limit int32,
) ([]flashcard, error) {
	sort := cur.Sort
	order := cur.Order
	position := cur.Position

	switch {
	case sort == createdAtSort && order == ascOrder && position != nil:
		return uc.db.ListFlashcardsByCreatedAtAscAfter(
			ctx,
			*position.CreatedAt,
			position.ID,
			limit,
		)

	case sort == createdAtSort && order == ascOrder:
		return uc.db.ListFlashcardsByCreatedAtAsc(ctx, limit)

	case sort == createdAtSort && order == descOrder && position != nil:
		return uc.db.ListFlashcardsByCreatedAtDescAfter(
			ctx,
			*position.CreatedAt,
			position.ID,
			limit,
		)
	case sort == createdAtSort && order == descOrder:
		return uc.db.ListFlashcardsByCreatedAtDesc(ctx, limit)
	default:
		return uc.db.ListFlashcardsByCreatedAtDesc(ctx, limit)
	}
}

func afterCursorFrom(sort sort, order order, last flashcard) (*cursor, error) {
	cur := &cursor{
		Sort:  sort,
		Order: order,
	}
	switch sort {
	case createdAtSort:
		cur.Position = &position{
			CreatedAt: &last.CreatedAt,
			ID:        last.ID,
		}
		return cur, nil
	default:
		return nil, fmt.Errorf("afterCursorFrom: unsupported sort %q", sort)
	}
}

type listFlashcardQueryResult struct {
	data  []flashcard
	after *cursor
}

type order string

const (
	ascOrder  order = "asc"
	descOrder order = "desc"
)

type sort string

const (
	createdAtSort sort = "created_at"
)

type cursor struct {
	Sort     sort      `json:"sort"`
	Order    order     `json:"order"`
	Position *position `json:"position"`
}

type position struct {
	ID        int32      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
