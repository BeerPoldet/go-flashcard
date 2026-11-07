package flashcardlist

import (
	"time"
)

type flashcard struct {
	ID        int32
	Title     string
	Content   string
	CreatedAt time.Time
}
