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

type IssueRequest struct {
	ProviderID  string
	DomainNames []string
	Parameters  map[string]any
}

type AvailableProvider struct {
	provider *Provider
}

func (a *AvailableProvider) ID() string {
	return (*a.provider).ID()
}

func (a *AvailableProvider) Name() string {
	return (*a.provider).Name()
}

func (a *AvailableProvider) DynamicFields() []*dynamic_fields.DynamicField {
	return (*a.provider).DynamicFields()
}

func (a *AvailableProvider) Priority() int {
	return (*a.provider).Priority()
}
