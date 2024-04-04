package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/r3d5un/modularmonolith/internal/logging"
	"github.com/r3d5un/modularmonolith/internal/peppol"
)

const (
	peppolBusinessCardPath string = "/api/v1/warehouse/peppolbusinesscards/"
)

type PeppolBusinessCard struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	CountryCode        string               `json:"countrycode"`
	LastUpdated        *time.Time           `json:"last_updated"`
	PeppolBusinessCard *peppol.BusinessCard `json:"peppol_business_card"`
}

type PeppolBusinessCardRequests struct {
	URL string `json:"url"`
}

func (r *PeppolBusinessCardRequests) Get(
	ctx context.Context,
	id string,
) (*PeppolBusinessCard, error) {
	logger := logging.LoggerFromContext(ctx)

	url := fmt.Sprintf("%s%s%s", r.URL, peppolBusinessCardPath, id)

	logger.Info("creating request", "url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("unable to create request", "error", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("unable to perform request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	var pbc PeppolBusinessCard

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("unable to read response body", "error", err)
		return nil, err
	}

	err = json.Unmarshal(body, &pbc)
	if err != nil {
		logger.Error("unable to unmarshal response body", "error", err)
		return nil, err
	}

	return &pbc, nil
}
