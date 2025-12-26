package certificate

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
)

type availableProviderResponse struct {
	ID            string                              `json:"id"`
	Name          string                              `json:"name"`
	Priority      int                                 `json:"priority"`
	DynamicFields []dynamicfield.DynamicFieldResponse `json:"dynamicFields"`
}

type certificateResponse struct {
	ID          uuid.UUID      `json:"id"`
	DomainNames []string       `json:"domainNames"`
	ProviderID  string         `json:"providerId"`
	IssuedAt    time.Time      `json:"issuedAt"`
	ValidUntil  time.Time      `json:"validUntil"`
	ValidFrom   time.Time      `json:"validFrom"`
	RenewAfter  *time.Time     `json:"renewAfter"`
	Parameters  map[string]any `json:"parameters"`
}

type issueCertificateRequest struct {
	ProviderID  string         `json:"providerId"`
	DomainNames []string       `json:"domainNames"`
	Parameters  map[string]any `json:"parameters"`
}

type issueCertificateResponse struct {
	Success       bool       `json:"success"`
	ErrorReason   *string    `json:"errorReason"`
	CertificateID *uuid.UUID `json:"certificateId"`
}

type renewCertificateResponse struct {
	Success     bool    `json:"success"`
	ErrorReason *string `json:"errorReason"`
}
