package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type certificateModel struct {
	bun.BaseModel `bun:"client_certificate"`

	ID                        uuid.UUID `bun:"id,pk"`
	Name                      string    `bun:"name"`
	Type                      string    `bun:"type"`
	ValidationMode            string    `bun:"validation_mode"`
	StaplingEnabled           bool      `bun:"stapling_enabled"`
	StaplingVerify            bool      `bun:"stapling_verify"`
	StaplingResponderURL      *string   `bun:"stapling_responder_url"`
	StaplingResponderFilePath *string   `bun:"stapling_responder_file_path"`
	CAPublicKey               *[]byte   `bun:"ca_public_key"`
	CAPrivateKey              *[]byte   `bun:"ca_private_key"`
	CASendToClients           *bool     `bun:"ca_send_to_clients"`
}

type clientModel struct {
	bun.BaseModel `bun:"client_certificate_item"`

	ID                  uuid.UUID  `bun:"id,pk"`
	ClientCertificateID uuid.UUID  `bun:"client_certificate_id"`
	DN                  *string    `bun:"dn"`
	PublicKey           *[]byte    `bun:"public_key"`
	PrivateKey          *[]byte    `bun:"private_key"`
	IssuedAt            *time.Time `bun:"issued_at"`
	ExpiresAt           *time.Time `bun:"expires_at"`
	Revoked             *bool      `bun:"revoked"`
}
