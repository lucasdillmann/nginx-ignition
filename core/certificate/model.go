package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
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

type AvailableProvider struct {
	ID            string
	Name          string
	Priority      int
	DynamicFields *[]*dynamic_fields.DynamicField
}

type IssueRequest struct {
	ProviderID  string
	DomainNames []*string
	Parameters  map[string]*any
}
