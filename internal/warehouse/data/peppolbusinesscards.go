package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/r3d5un/modularmonolith/internal/logging"
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
	DB      *sql.DB
	Timeout *time.Duration
}

func (m *PeppolBusinessCardModel) Get(ctx context.Context, id string) (*PeppolBusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	stmt := `SELECT id, name, countrycode, last_updated, business_card
FROM peppol_business_cards
WHERE id = $1;`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var pbc PeppolBusinessCard
	var jpbc []byte

	logger.LogAttrs(
		qCtx, slog.LevelInfo, "querying peppol business card",
		slog.String("string", stmt), slog.String("id", id),
	)
	row := m.DB.QueryRowContext(qCtx, stmt, id)
	err := row.Scan(&pbc.ID, &pbc.Name, &pbc.CountryCode, &pbc.LastUpdated, jpbc)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.LogAttrs(qCtx, slog.LevelInfo, "no record found")
			return nil, ErrNoRecord
		default:
			logger.Error("unable to query business card", "error", err)
			return nil, err
		}
	}
	logger.Info("data retrieved")

	err = json.Unmarshal(jpbc, &pbc.PeppolBusinessCard)
	if err != nil {
		logger.Error("error unmarhaling peppol business card", "error", err)
		return nil, err
	}

	return &pbc, nil
}

func (m *PeppolBusinessCardModel) Insert(
	ctx context.Context,
	pbc *PeppolBusinessCard,
) (*PeppolBusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	stmt := `INSERT INTO peppol_business_cards (id, name, countrycode, last_updated, business_card)
VALUES ($1, $2, $3, NOW(), $5)
RETURNING id, name, countrycode, last_updated, business_card;`

	pbcBytes, err := json.Marshal(pbc.PeppolBusinessCard)
	if err != nil {
		logger.Error("error marshaling peppol_business_card", "error", err)
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var jpbc []byte

	logger.Info("inserting peppol business card", "query", stmt, "id", pbc.ID)
	row := m.DB.QueryRowContext(qCtx, stmt, pbc.ID, pbc.Name, pbc.CountryCode, pbcBytes)
	err = row.Scan(&pbc.ID, &pbc.Name, &pbc.CountryCode, &pbc.LastUpdated, &jpbc)
	if err != nil {
		logger.Error("error inserting peppol business card", "error", err)
		return nil, err
	}
	logger.InfoContext(ctx, "peppol business card inserted", "id", pbc.ID)

	err = json.Unmarshal(jpbc, &pbc.PeppolBusinessCard)
	if err != nil {
		logger.ErrorContext(ctx, "error unmarshaling peppol_business_card", "error", err)
		return nil, err
	}

	return pbc, nil
}

func (m *PeppolBusinessCardModel) Upsert(
	ctx context.Context,
	pbc *PeppolBusinessCard,
) (*PeppolBusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	stmt := `INSERT INTO peppol_business_cards (id, name, countrycode, last_updated, business_card)
VALUES ($1, $2, $3, NOW(), $4)
ON CONFLICT (id) DO UPDATE
SET
    name = EXCLUDED.name, countrycode = EXCLUDED.countrycode,
    last_updated = NOW(), business_card = EXCLUDED.business_card
RETURNING id, name, countrycode, last_updated, business_card;`

	pbcBytes, err := json.Marshal(pbc.PeppolBusinessCard)
	if err != nil {
		logger.Error("error marshaling peppol_business_card", "error", err)
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var jpbc []byte

	logger.Info("inserting peppol business card", "query", stmt, "id", pbc.ID)
	row := m.DB.QueryRowContext(qCtx, stmt, pbc.ID, pbc.Name, pbc.CountryCode, pbcBytes)
	err = row.Scan(&pbc.ID, &pbc.Name, &pbc.CountryCode, &pbc.LastUpdated, &jpbc)
	if err != nil {
		logger.Error("error inserting peppol business card", "error", err)
		return nil, err
	}
	logger.InfoContext(ctx, "peppol business card inserted", "id", pbc.ID)

	err = json.Unmarshal(jpbc, &pbc.PeppolBusinessCard)
	if err != nil {
		logger.ErrorContext(ctx, "error unmarshaling peppol_business_card", "error", err)
		return nil, err
	}

	return pbc, nil
}

func (m *PeppolBusinessCardModel) Delete(
	ctx context.Context,
	id string,
) (*PeppolBusinessCard, error) {
	return nil, nil
}

func (m *PeppolBusinessCardModel) Update(
	ctx context.Context,
	pbc *PeppolBusinessCard,
) (*PeppolBusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	stmt := `UPDATE peppol_business_card
SET id = $1, name = $2, countrycode = $3, last_updated = NOW(), business_card = $4
WHERE id = $1
RETURNING id, name, countrycode, last_updated, business_card;
    `

	jsonBc, err := json.Marshal(pbc.PeppolBusinessCard)
	if err != nil {
		return nil, err
	}

	args := []any{
		pbc.ID,
		pbc.Name,
		pbc.CountryCode,
		jsonBc,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var bcRaw []byte

	logger.InfoContext(rCtx, "updating business card", "query", stmt, "args", args)
	err = m.DB.QueryRowContext(ctx, stmt, args...).Scan(
		&pbc.ID, &pbc.Name, &pbc.CountryCode, &pbc.LastUpdated, &bcRaw,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.InfoContext(ctx, "no records found", "query", stmt, "args", args)
			return nil, ErrNoRecord
		default:
			logger.ErrorContext(
				ctx, "unable to update business card",
				"query", stmt, "args", args, "error", err,
			)
			return nil, err
		}
	}
	logger.InfoContext(
		rCtx,
		"updated brreg business card",
		"business_card",
		pbc.PeppolBusinessCard,
	)

	err = json.Unmarshal(bcRaw, &pbc.PeppolBusinessCard)
	if err != nil {
		logger.ErrorContext(ctx, "error unmarshaling business card", "error", err)
		return nil, err
	}

	return pbc, nil
}
