package letsencrypt

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

type validationRules struct {
	dynamicFields []dynamicfields.DynamicField
}

func (r validationRules) DynamicFields() []dynamicfields.DynamicField {
	return r.dynamicFields
}

func (r validationRules) Validate(
	ctx context.Context,
	request *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	output := make([]validation.ConsistencyViolation, 0)

	termsOfServiceAccepted, casted := request.Parameters[termsOfServiceFieldID].(bool)
	if !casted || !termsOfServiceAccepted {
		output = append(output, validation.ConsistencyViolation{
			Path:    fmt.Sprintf("parameters.%s", termsOfServiceFieldID),
			Message: i18n.M(ctx, i18n.K.CertificateLetsencryptTosRequired),
		})
	}

	_, dnsViolations := resolveCustomDNSServers(ctx, request.Parameters)
	output = append(output, dnsViolations...)

	return output
}

func resolveCustomDNSServers(
	ctx context.Context,
	parameters map[string]any,
) ([]string, []validation.ConsistencyViolation) {
	bypass, _ := parameters[bypassDNSPropagationChecks].(bool)
	if bypass {
		return nil, nil
	}

	override, _ := parameters[overrideDNSServersFieldID].(bool)
	if !override {
		return nil, nil
	}

	raw, _ := parameters[customDNSServersFieldID].(string)
	if raw == "" {
		return nil, nil
	}

	servers := make([]string, 0)
	violations := make([]validation.ConsistencyViolation, 0)
	path := fmt.Sprintf("parameters.%s", customDNSServersFieldID)

	for _, token := range strings.Fields(raw) {
		if !isValidDNSServer(token) {
			violations = append(violations, validation.ConsistencyViolation{
				Path: path,
				Message: i18n.M(ctx, i18n.K.CertificateLetsencryptInvalidDnsServers).
					V("value", token),
			})

			continue
		}

		servers = append(servers, token)
	}

	if len(violations) > 0 {
		return nil, violations
	}

	if len(servers) == 0 {
		return nil, []validation.ConsistencyViolation{{
			Path:    path,
			Message: i18n.M(ctx, i18n.K.CertificateLetsencryptInvalidDnsServers).V("value", ""),
		}}
	}

	return servers, nil
}
