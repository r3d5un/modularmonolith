package data

import (
	"database/sql"
	"time"
)

type Models struct {
	PeppolBusinessCards PeppolBusinessCardModel
}

func NewModels(db *sql.DB, timeout *time.Duration) Models {
	return Models{
		PeppolBusinessCards: PeppolBusinessCardModel{DB: db, Timeout: timeout},
	}
}
