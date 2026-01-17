package certificate

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Certificate struct {
	IssuedAt           time.Time
	ValidUntil         time.Time
	ValidFrom          time.Time
	RenewAfter         *time.Time
	Parameters         map[string]any
	Metadata           *string
	ProviderID         string
	PrivateKey         string
	PublicKey          string
	DomainNames        []string
	CertificationChain []string
	ID                 uuid.UUID
}

type IssueRequest struct {
	Parameters  map[string]any
	ProviderID  string
	DomainNames []string
}

type AutoRenewSettings struct {
	IntervalUnit      string
	IntervalUnitCount int
	Enabled           bool
}

type AvailableProvider struct {
	provider Provider
}

func (a *AvailableProvider) ID() string {
	return a.provider.ID()
}

func (a *AvailableProvider) Name(ctx context.Context) *i18n.Message {
	return a.provider.Name(ctx)
}

func (a *AvailableProvider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return a.provider.DynamicFields(ctx)
}

func (a *AvailableProvider) Priority() int {
	return a.provider.Priority()
}
