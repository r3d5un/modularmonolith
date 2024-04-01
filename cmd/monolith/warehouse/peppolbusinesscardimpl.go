package warehouse

import (
	"context"
	"errors"
	"log/slog"

	"github.com/r3d5un/modularmonolith/internal/logging"
	"github.com/r3d5un/modularmonolith/internal/peppol"
	"github.com/r3d5un/modularmonolith/internal/warehouse/data"
)

func (m *Module) GetPeppolBusinessCard(
	ctx context.Context,
	id string,
) (*peppol.BusinessCard, error) {
	logger := enrichContextEmbeddedLogger(ctx)

	logger.Info("getting data", "id", id)
	pbcRecord, err := m.models.PeppolBusinessCards.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			m.log.Info("peppol business card not found")
			return nil, err
		default:
			m.log.Error("unable to get peppol business card", "error", err)
			return nil, err
		}
	}

	return pbcRecord.PeppolBusinessCard, nil
}

func enrichContextEmbeddedLogger(ctx context.Context) *slog.Logger {
	contextLogger := logging.LoggerFromContext(ctx)
	enrichedLogger := contextLogger.With(slog.Group("module", slog.String("name", "warehouse")))
	return enrichedLogger
}
