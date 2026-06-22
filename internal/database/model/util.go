package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func fromString(db pgtype.Text) string {
	if db.Valid {
		return db.String
	}

	return ""
}

func fromInt(db pgtype.Int4) int {
	if db.Valid {
		return int(db.Int32)
	}

	return 0
}

func fromTime(db pgtype.Timestamptz) time.Time {
	if db.Valid {
		return db.Time
	}

	return time.Time{}
}
