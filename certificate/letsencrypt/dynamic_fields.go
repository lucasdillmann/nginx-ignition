package letsencrypt

import (
	"sort"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

var (
	termsOfService = dynamicfields.DynamicField{
		ID:          "acceptTheTermsOfService",
		Priority:    99,
		Description: "Terms of service",
		HelpText: ptr.Of(
			"I agree to the Let's Encrypt terms of service available at theirs website",
		),
		Required:     true,
		DefaultValue: false,
		Type:         dynamicfields.BooleanType,
	}

	emailAddress = dynamicfields.DynamicField{
		ID:          "emailAddress",
		Priority:    0,
		Description: "E-mail address",
		Required:    true,
		Type:        dynamicfields.EmailType,
	}

	dnsProvider = dynamicfields.DynamicField{
		ID:          "challengeDnsProvider",
		Priority:    1,
		Description: "DNS provider",
		Required:    true,
		Type:        dynamicfields.EnumType,
	}
)

func resolveDynamicFields() []dynamicfields.DynamicField {
	output := make([]dynamicfields.DynamicField, 0, len(providers))

	output = append(output, termsOfService, emailAddress, dnsProvider)
	providerOptions := make([]dynamicfields.EnumOption, 0, len(providers))

	for _, provider := range providers {
		output = append(output, provider.DynamicFields()...)

		providerOptions = append(providerOptions, dynamicfields.EnumOption{
			ID:          provider.ID(),
			Description: provider.Name(),
		})
	}

	sort.Slice(providerOptions, func(leftIndex, rightIndex int) bool {
		leftValue := strings.ToUpper(providerOptions[leftIndex].Description)
		rightValue := strings.ToUpper(providerOptions[rightIndex].Description)

		return leftValue < rightValue
	})

	dnsProvider.EnumOptions = providerOptions

	return output
}
