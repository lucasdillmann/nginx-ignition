package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"github.com/google/uuid"
	"time"
)

type availableProviderResponse struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Priority      int                     `json:"priority"`
	DynamicFields []*dynamicFieldResponse `json:"dynamicFields"`
}

type dynamicFieldResponse struct {
	Name        string              `json:"name"`
	Type        dynamic_fields.Type `json:"type"`
	Description string              `json:"description"`
}

type certificateResponse struct {
	ID          uuid.UUID              `json:"id"`
	DomainNames []string               `json:"domainNames"`
	ProviderID  string                 `json:"providerId"`
	IssuedAt    time.Time              `json:"issuedAt"`
	ValidUntil  time.Time              `json:"validUntil"`
	ValidFrom   time.Time              `json:"validFrom"`
	RenewAfter  *time.Time             `json:"renewAfter"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type issueCertificateRequest struct {
	ProviderID  string                  `json:"providerId" validation:"required"`
	DomainNames []*string               `json:"domainNames" validation:"required,nonempty"`
	Parameters  map[string]*interface{} `json:"parameters"`
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
