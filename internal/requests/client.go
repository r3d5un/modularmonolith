package requests

import (
	"context"

	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/logging"
	"github.com/r3d5un/modularmonolith/internal/peppol"
)

type Client struct {
	Requests Requests
}

func NewClient(config config.RequestConfiguration) *Client {
	c := Client{Requests: NewRequests(config)}

	return &c
}

func (c *Client) GetPeppolBusinessCard(
	ctx context.Context,
	id string,
) (*peppol.BusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	logger.Info("getting data", "id", id)
	pbc, err := c.Requests.PeppolBusinessCards.Get(ctx, id)
	if err != nil {
		logger.Error("unable to retrieve peppol business card", "error", err)
		return nil, err
	}

	return pbc.PeppolBusinessCard, nil
}
