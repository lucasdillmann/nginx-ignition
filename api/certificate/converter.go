package certificate

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamic_field"
	"dillmann.com.br/nginx-ignition/core/certificate"
)

func toAvailableProviderResponse(input []*certificate.AvailableProvider) []*availableProviderResponse {
	var responses []*availableProviderResponse
	for _, provider := range input {
		responses = append(responses, &availableProviderResponse{
			ID:            provider.ID(),
			Name:          provider.Name(),
			Priority:      provider.Priority(),
			DynamicFields: dynamic_field.ToResponse(provider.DynamicFields()),
		})
	}
	return responses
}

func toIssueCertificateResponse(certificate *certificate.Certificate, err error) *issueCertificateResponse {
	var errorReason *string
	if err != nil {
		errorStr := err.Error()
		errorReason = &errorStr
	}

	var certificateId *uuid.UUID
	if certificate != nil {
		certificateId = &certificate.ID
	}

	return &issueCertificateResponse{
		Success:       err == nil,
		ErrorReason:   errorReason,
		CertificateID: certificateId,
	}
}

func toRenewCertificateResponse(err error) *renewCertificateResponse {
	var errorReason *string
	if err != nil {
		errorStr := err.Error()
		errorReason = &errorStr
	}

	return &renewCertificateResponse{
		Success:     err == nil,
		ErrorReason: errorReason,
	}
}

func toCertificateResponse(input *certificate.Certificate) *certificateResponse {
	return &certificateResponse{
		ID:          input.ID,
		DomainNames: input.DomainNames,
		ProviderID:  input.ProviderID,
		IssuedAt:    input.IssuedAt,
		ValidUntil:  input.ValidUntil,
		ValidFrom:   input.ValidFrom,
		RenewAfter:  input.RenewAfter,
		Parameters:  input.Parameters,
	}
}

func toIssueCertificateRequest(input *issueCertificateRequest) *certificate.IssueRequest {
	return &certificate.IssueRequest{
		ProviderID:  input.ProviderID,
		DomainNames: dereference(input.DomainNames),
		Parameters:  input.Parameters,
	}
}

func dereference(input []*string) []string {
	var output []string
	for _, str := range input {
		if str != nil {
			output = append(output, *str)
		}
	}

	return output
}
