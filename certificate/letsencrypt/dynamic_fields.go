package letsencrypt

import (
	"sort"

	"github.com/aws/smithy-go/ptr"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

var (
	termsOfService = dynamic_fields.DynamicField{
		ID:          "acceptTheTermsOfService",
		Priority:    99,
		Description: "Terms of service",
		HelpText:    ptr.String("I agree to the Let's Encrypt terms of service available at theirs website"),
		Required:    true,
		Type:        dynamic_fields.BooleanType,
	}

	emailAddress = dynamic_fields.DynamicField{
		ID:          "emailAddress",
		Priority:    0,
		Description: "E-mail address",
		Required:    true,
		Type:        dynamic_fields.EmailType,
	}

	dnsProvider = dynamic_fields.DynamicField{
		ID:          "challengeDnsProvider",
		Priority:    1,
		Description: "DNS provider (for the DNS challenge)",
		Required:    true,
		Type:        dynamic_fields.EnumType,
	}
)

func resolveDynamicFields() []*dynamic_fields.DynamicField {
	output := make([]*dynamic_fields.DynamicField, 0, len(providers))

	output = append(output, &termsOfService, &emailAddress, &dnsProvider)
	providerOptions := make([]*dynamic_fields.EnumOption, 0, len(providers))

	for _, provider := range providers {
		output = append(output, provider.DynamicFields()...)

		providerOptions = append(providerOptions, &dynamic_fields.EnumOption{
			ID:          provider.ID(),
			Description: provider.Name(),
		})
	}

	sort.Slice(providerOptions, func(leftIndex, rightIndex int) bool {
		return providerOptions[leftIndex].Description < providerOptions[rightIndex].Description
	})

	dnsProvider.EnumOptions = &providerOptions

	return output
}
