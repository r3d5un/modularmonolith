package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/r3d5un/modularmonolith/internal/peppol"
)

type PeppolBusinessCard struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	CountryCode        string               `json:"countrycode"`
	LastUpdated        *time.Time           `json:"last_updated"`
	PeppolBusinessCard *peppol.BusinessCard `json:"peppol_business_card"`
}

type PeppolBusinessCardModel struct {
	DB *sql.DB
}

func (m *PeppolBusinessCardModel) Get(ctx context.Context, id string) (*PeppolBusinessCard, error) {
	return nil, nil
}

func (m *PeppolBusinessCardModel) Insert(
	ctx context.Context,
	pbc PeppolBusinessCard,
) (*PeppolBusinessCard, error) {
	return nil, nil
}

func (m *PeppolBusinessCardModel) Upsert(
	ctx context.Context,
	bc *PeppolBusinessCard,
) (*PeppolBusinessCard, error) {
	return nil, nil
}

func (m *PeppolBusinessCardModel) Delete(
	ctx context.Context,
	id string,
) (*PeppolBusinessCard, error) {
	return nil, nil
}
