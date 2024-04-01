package warehouse

import (
	"errors"
	"net/http"

	"github.com/r3d5un/modularmonolith/internal/peppol"
	"github.com/r3d5un/modularmonolith/internal/utils/httputils"
	"github.com/r3d5un/modularmonolith/internal/warehouse/data"
)

type PeoolBusinessCardResponse struct {
	Data data.PeppolBusinessCard
}

func (m *Module) getPeppolBusinessCardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	m.log.Info("parsing ID")
	s := r.PathValue("id")
	if s == "" {
		m.log.Info("parameter empty", "key", "id")
		httputils.BadRequestResponse(w, r, "id parameter emtpy")
		return
	}

	m.log.Info("getting data", "id", s)
	pbc, err := m.models.PeppolBusinessCards.Get(ctx, s)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			m.log.Info("peppol business card not found")
			httputils.NotFoundResponse(w, r)
			return
		default:
			m.log.Error("unable to get peppol business card", "error", err)
			httputils.ServerErrorResponse(w, r, err)
			return
		}
	}

	m.log.Info("returning peppol business card")
	err = httputils.WriteJSON(
		w, http.StatusOK, PeoolBusinessCardResponse{Data: *pbc}, nil,
	)
	if err != nil {
		m.log.Error("error writing response", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *Module) deletePeppolBusinessCardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	m.log.Info("parsing ID")
	s := r.PathValue("id")
	if s == "" {
		m.log.Info("parameter empty", "key", "id")
		httputils.BadRequestResponse(w, r, "id parameter emtpy")
		return
	}

	m.log.Info("getting data", "id", s)
	pbc, err := m.models.PeppolBusinessCards.Delete(ctx, s)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			m.log.Info("peppol business card not found")
			httputils.NotFoundResponse(w, r)
			return
		default:
			m.log.Error("unable to delete peppol business card", "error", err)
			httputils.ServerErrorResponse(w, r, err)
			return
		}
	}

	m.log.Info("returning peppol business card")
	err = httputils.WriteJSON(
		w, http.StatusOK, PeoolBusinessCardResponse{Data: *pbc}, nil,
	)
	if err != nil {
		m.log.Error("error writing response", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *Module) putPeppolBusinessCardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var pbcBody *peppol.BusinessCard

	m.log.Info("parsing business card from request body")
	err := httputils.ReadJSON(r, pbcBody)
	if err != nil {
		m.log.Error("unable to read peppol business card", "error", err)
		httputils.BadRequestResponse(w, r, "unable to read peppol business card")
		return
	}

	pbc := data.PeppolBusinessCard{
		ID:                 pbcBody.Participant.Value,
		Name:               pbcBody.Entity.Name.Name,
		CountryCode:        pbcBody.Entity.CountryCode,
		PeppolBusinessCard: pbcBody,
	}

	m.log.Info("upserting peppol business card")
	upsertedPBC, err := m.models.PeppolBusinessCards.Upsert(ctx, &pbc)
	if err != nil {
		m.log.Error("unable to upsert business card", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}

	m.log.Info("returning peppol business card")
	err = httputils.WriteJSON(
		w, http.StatusOK, PeoolBusinessCardResponse{Data: *upsertedPBC}, nil,
	)
	if err != nil {
		m.log.Error("error writing response", "error", err)
		httputils.ServerErrorResponse(w, r, err)
		return
	}
}
