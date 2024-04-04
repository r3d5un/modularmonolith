package requests

import (
	"github.com/r3d5un/modularmonolith/internal/config"
)

type Requests struct {
	PeppolBusinessCards PeppolBusinessCardRequests
}

func NewRequests(config config.RequestConfiguration) Requests {
	return Requests{
		PeppolBusinessCards: PeppolBusinessCardRequests{config.URL},
	}
}
