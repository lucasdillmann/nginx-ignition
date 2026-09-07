package acme

import (
	"context"
	"sort"
	"strings"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

//nolint:gosec
const (
	termsOfServiceFieldID       = "acceptTheTermsOfService"
	emailAddressFieldID         = "emailAddress"
	dnsProviderFieldID          = "challengeDnsProvider"
	bypassDNSPropagationChecks  = "bypassDnsPropagationChecks"
	certificateAuthorityFieldID = "certificateAuthority"

	LetsEncryptProduction = "LETSENCRYPT_PRODUCTION"
	LetsEncryptSandbox    = "LETSENCRYPT_SANDBOX"
	ActalisProduction     = "ACTALIS_PRODUCTION"
	ZeroSSLProduction     = "ZEROSSL_PRODUCTION"
)

func mainDynamicFields(ctx context.Context) ([]dynamicfields.DynamicField, int) {
	certificateAuthorityField := dynamicfields.DynamicField{
		ID:           certificateAuthorityFieldID,
		Priority:     1,
		Description:  i18n.M(ctx, i18n.K.CertificateAcmeDnsProvider),
		Required:     true,
		Type:         dynamicfields.EnumType,
		DefaultValue: LetsEncryptProduction,
		EnumOptions: []dynamicfields.EnumOption{
			{
				ID:          LetsEncryptProduction,
				Description: i18n.M(ctx, i18n.K.CertificateAcmeLetsencryptProduction),
			},
			{
				ID:          LetsEncryptSandbox,
				Description: i18n.M(ctx, i18n.K.CertificateAcmeLetsencryptSandbox),
			},
			{
				ID:          ActalisProduction,
				Description: i18n.M(ctx, i18n.K.CertificateAcmeActalisProduction),
			},
			{
				ID:          ZeroSSLProduction,
				Description: i18n.M(ctx, i18n.K.CertificateAcmeZerosslProduction),
			},
		},
	}

	dnsField := dynamicfields.DynamicField{
		ID:          dnsProviderFieldID,
		Priority:    2,
		Description: i18n.M(ctx, i18n.K.CertificateAcmeDnsProvider),
		Required:    true,
		Type:        dynamicfields.EnumType,
	}

	bypassDNSPropagationChecksField := dynamicfields.DynamicField{
		ID:           bypassDNSPropagationChecks,
		Priority:     98,
		Description:  i18n.M(ctx, i18n.K.CertificateAcmeBypassDnsPropagationChecks),
		Required:     false,
		DefaultValue: false,
		Type:         dynamicfields.BooleanType,
	}

	tosField := dynamicfields.DynamicField{
		ID:           termsOfServiceFieldID,
		Priority:     99,
		Description:  i18n.M(ctx, i18n.K.CertificateAcmeTos),
		HelpText:     i18n.M(ctx, i18n.K.CertificateAcmeTosHelp),
		Required:     true,
		DefaultValue: false,
		Type:         dynamicfields.BooleanType,
	}

	emailField := dynamicfields.DynamicField{
		ID:          emailAddressFieldID,
		Priority:    0,
		Description: i18n.M(ctx, i18n.K.CertificateAcmeEmail),
		Required:    true,
		Type:        dynamicfields.EmailType,
	}

	return []dynamicfields.DynamicField{
		certificateAuthorityField,
		dnsField,
		emailField,
		bypassDNSPropagationChecksField,
		tosField,
	}, 1
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
