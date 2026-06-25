package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func toString(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func toInt(value int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(value), Valid: value != 0}
}

func toBool(value *bool) pgtype.Bool {
	v := false
	if value != nil {
		v = *value
	}

	return pgtype.Bool{Bool: v, Valid: value != nil}
}

func toTime(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: !value.IsZero()}
}
