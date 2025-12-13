package server

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID             uuid.UUID
	Type           Type
	CA             CACertificate
	Clients        []ClientCertificate
	ValidationMode ValidationMode
	Stapling       Stapling
}

type Type string

const (
	ManagedType  Type = "MANAGED"
	ExternalType Type = "EXTERNAL"
)

type CACertificate struct {
	DN            *string
	PublicKey     *string
	PrivateKey    *string
	SendToClients *bool
}

type ClientCertificate struct {
	DN         *string
	PublicKey  *string
	PrivateKey *string
	IssuedAt   *time.Time
}

type Stapling struct {
	Enabled           bool
	Verify            bool
	ResponderURL      *string
	ResponderFilePath *string
}

type ValidationMode string

const (
	OnValidationMode           ValidationMode = "ON"
	OffValidationMode          ValidationMode = "OFF"
	OptionalValidationMode     ValidationMode = "OPTIONAL"
	OptionalNoCAValidationMode ValidationMode = "OPTIONAL_NO_CA"
)
