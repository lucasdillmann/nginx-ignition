package client

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID             uuid.UUID
	Name           string
	Type           Type
	CA             *CA
	Clients        []*Client
	ValidationMode ValidationMode
	Stapling       Stapling
}

type Type string

const (
	SimpleType            Type = "SIMPLE"
	IgnitionManagedType   Type = "IGNITION_MANAGED"
	ExternallyManagedType Type = "EXTERNALLY_MANAGED"
)

type CA struct {
	PublicKey     *[]byte
	PrivateKey    *[]byte
	SendToClients *bool
}

type Client struct {
	ID         uuid.UUID
	DN         *string
	PublicKey  *[]byte
	PrivateKey *[]byte
	IssuedAt   *time.Time
	ExpiresAt  *time.Time
	Revoked    *bool
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

type CreateRequest struct {
	Name         string
	Type         Type
	CARequest    *CertificateRequest
	CAPublicKey  *[]byte
	CAPrivateKey *[]byte
}

type UpdateRequest struct {
	Name           string
	ValidationMode ValidationMode
	Stapling       Stapling
}

type ReplaceCARequest struct {
	CARequest    *CertificateRequest
	CAPublicKey  *[]byte
	CAPrivateKey *[]byte
}

type CreateClientRequest struct {
	DN      *string
	Request *CertificateRequest
}

type UpdateClientRequest struct {
	ID      uuid.UUID
	DN      *string
	Revoked *bool
}
