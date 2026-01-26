package letsencrypt

import (
	"context"
	"sort"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	termsOfServiceFieldID = "acceptTheTermsOfService"
	emailAddressFieldID   = "emailAddress"
	dnsProviderFieldID    = "challengeDnsProvider"
)

func mainDynamicFields(ctx context.Context) ([]dynamicfields.DynamicField, int) {
	dnsField := dynamicfields.DynamicField{
		ID:          dnsProviderFieldID,
		Priority:    1,
		Description: i18n.M(ctx, i18n.K.CertificateLetsencryptDnsProvider),
		Required:    true,
		Type:        dynamicfields.EnumType,
	}

	tosField := dynamicfields.DynamicField{
		ID:           termsOfServiceFieldID,
		Priority:     99,
		Description:  i18n.M(ctx, i18n.K.CertificateLetsencryptTos),
		HelpText:     i18n.M(ctx, i18n.K.CertificateLetsencryptTosHelp),
		Required:     true,
		DefaultValue: false,
		Type:         dynamicfields.BooleanType,
	}

	emailField := dynamicfields.DynamicField{
		ID:          emailAddressFieldID,
		Priority:    0,
		Description: i18n.M(ctx, i18n.K.CertificateLetsencryptEmail),
		Required:    true,
		Type:        dynamicfields.EmailType,
	}

	return []dynamicfields.DynamicField{dnsField, tosField, emailField}, 0
}

func resolveDynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	mainFields, dnsProviderField := mainDynamicFields(ctx)
	output := make([]dynamicfields.DynamicField, 0, 3+len(providers))
	output = append(output, mainFields...)
	providerOptions := make([]dynamicfields.EnumOption, 0, len(providers))

	for _, provider := range providers {
		output = append(output, provider.DynamicFields(ctx)...)

		providerOptions = append(providerOptions, dynamicfields.EnumOption{
			ID:          provider.ID(),
			Description: provider.Name(ctx),
		})
	}

	sort.Slice(providerOptions, func(leftIndex, rightIndex int) bool {
		leftValue := strings.ToUpper(providerOptions[leftIndex].Description.String())
		rightValue := strings.ToUpper(providerOptions[rightIndex].Description.String())

		return leftValue < rightValue
	})

	output[dnsProviderField].EnumOptions = providerOptions

	return output
}
