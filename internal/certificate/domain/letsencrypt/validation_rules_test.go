package letsencrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_validationRules_Validate_customDNSServers(t *testing.T) {
	tests := []struct {
		parameters         map[string]any
		name               string
		expectedViolations int
	}{
		{
			name: "valid servers",
			parameters: map[string]any{
				termsOfServiceFieldID:     true,
				overrideDNSServersFieldID: true,
				customDNSServersFieldID:   "1.1.1.1\n8.8.8.8:53",
			},
		},
		{
			name: "invalid server",
			parameters: map[string]any{
				termsOfServiceFieldID:     true,
				overrideDNSServersFieldID: true,
				customDNSServersFieldID:   "example.com",
			},
			expectedViolations: 1,
		},
		{
			name: "whitespace-only servers",
			parameters: map[string]any{
				termsOfServiceFieldID:     true,
				overrideDNSServersFieldID: true,
				customDNSServersFieldID:   "   \n\t",
			},
			expectedViolations: 1,
		},
		{
			name: "empty servers",
			parameters: map[string]any{
				termsOfServiceFieldID:     true,
				overrideDNSServersFieldID: true,
				customDNSServersFieldID:   "",
			},
		},
		{
			name: "override disabled",
			parameters: map[string]any{
				termsOfServiceFieldID:     true,
				overrideDNSServersFieldID: false,
				customDNSServersFieldID:   "example.com",
			},
		},
		{
			name: "propagation checks skipped",
			parameters: map[string]any{
				termsOfServiceFieldID:      true,
				bypassDNSPropagationChecks: true,
				overrideDNSServersFieldID:  true,
				customDNSServersFieldID:    "example.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules := validationRules{}
			violations := rules.Validate(t.Context(), &certificate.IssueRequest{
				Parameters: test.parameters,
			})

			assert.Len(t, violations, test.expectedViolations)

			if test.expectedViolations > 0 {
				assert.Equal(t, "parameters.customDnsServers", violations[0].Path)
			}
		})
	}
}
