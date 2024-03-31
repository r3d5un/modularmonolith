package data

import (
	"database/sql"
)

type Models struct {
	PeppolBusinessCards PeppolBusinessCardModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		PeppolBusinessCards: PeppolBusinessCardModel{DB: db},
	}
}
