-- name: FlashcardList :many
SELECT * FROM flashcards LIMIT $1;
