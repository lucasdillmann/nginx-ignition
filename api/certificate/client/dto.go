package client

import (
	"time"

	"github.com/google/uuid"
)

type certificateResponseDto struct {
	ID              uuid.UUID            `json:"id"`
	Name            string               `json:"name"`
	Type            string               `json:"type"`
	ValidationMode  string               `json:"validationMode"`
	SendCAToClients *bool                `json:"sendCaToClients,omitempty"`
	Clients         *[]clientResponseDto `json:"clients,omitempty"`
	Stapling        *staplingResponseDto `json:"stapling"`
}

type staplingResponseDto struct {
	Enabled           bool    `json:"enabled"`
	Verify            bool    `json:"verify"`
	ResponderURL      *string `json:"responderUrl"`
	ResponderFilePath *string `json:"responderFilePath"`
}

type clientResponseDto struct {
	ID        uuid.UUID  `json:"id"`
	DN        *string    `json:"dn,omitempty"`
	IssuedAt  *time.Time `json:"issuedAt,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Revoked   *bool      `json:"revoked,omitempty"`
}
