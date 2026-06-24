package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func toInt(value int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(value), Valid: value != 0}
}
