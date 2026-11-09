-- name: FlashcardListByCreatedAtAsc :many
SELECT * FROM flashcards
ORDER BY created_at ASC, id ASC
LIMIT @lim;

-- name: FlashcardListByCreatedAtDesc :many
SELECT * FROM flashcards
ORDER BY created_at DESC, id DESC
LIMIT @lim;

-- name: FlashcardListByCreatedAtAscAfter :many
SELECT * FROM flashcards
WHERE (created_at, id) > (@created_at::timestamp, @id::int)
ORDER BY created_at ASC, id ASC
LIMIT @lim;

-- name: FlashcardListByCreatedAtDescAfter :many
SELECT * FROM flashcards
WHERE (created_at, id) < (@created_at::timestamp, @id::int)
ORDER BY created_at DESC, id DESC
LIMIT @lim;
