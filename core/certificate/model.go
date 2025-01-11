package certificate

import (
	"github.com/google/uuid"
	"time"
)

type Certificate struct {
	ID                 uuid.UUID
	DomainNames        []string
	ProviderID         string
	IssuedAt           time.Time
	ValidUntil         time.Time
	ValidFrom          time.Time
	RenewAfter         *time.Time
	PrivateKey         string
	PublicKey          string
	CertificationChain []string
	Parameters         map[string]interface{}
	Metadata           *string
}
