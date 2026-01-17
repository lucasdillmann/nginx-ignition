package certificate

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type availableProviderResponse struct {
	ID            string                  `json:"id"`
	Name          *i18n.Message           `json:"name"`
	DynamicFields []dynamicfield.Response `json:"dynamicFields"`
	Priority      int                     `json:"priority"`
}

type certificateResponse struct {
	IssuedAt    time.Time      `json:"issuedAt"`
	ValidUntil  time.Time      `json:"validUntil"`
	ValidFrom   time.Time      `json:"validFrom"`
	RenewAfter  *time.Time     `json:"renewAfter"`
	Parameters  map[string]any `json:"parameters"`
	ProviderID  string         `json:"providerId"`
	DomainNames []string       `json:"domainNames"`
	ID          uuid.UUID      `json:"id"`
}

type issueCertificateRequest struct {
	Parameters  map[string]any `json:"parameters"`
	ProviderID  string         `json:"providerId"`
	DomainNames []string       `json:"domainNames"`
}

type issueCertificateResponse struct {
	ErrorReason   *string    `json:"errorReason"`
	CertificateID *uuid.UUID `json:"certificateId"`
	Success       bool       `json:"success"`
}

type renewCertificateResponse struct {
	ErrorReason *string `json:"errorReason"`
	Success     bool    `json:"success"`
}
