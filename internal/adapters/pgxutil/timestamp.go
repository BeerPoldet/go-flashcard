package pgxutil

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}
