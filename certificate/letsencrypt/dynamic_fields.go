package letsencrypt

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

var (
	awsRoute53Id  = "AWS_ROUTE53"
	azureId       = "AZURE"
	cloudflareId  = "CLOUDFLARE"
	googleCloudId = "GOOGLE_CLOUD"

	termsOfService = dynamic_fields.DynamicField{
		ID:          "acceptTheTermsOfService",
		Priority:    99,
		Description: "Terms of service",
		HelpText:    stringPtr("I agree to the Let's Encrypt terms of service available at theirs website"),
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
		EnumOptions: &[]*dynamic_fields.EnumOption{
			{ID: awsRoute53Id, Description: "AWS Route53"},
			{ID: azureId, Description: "Azure DNS"},
			{ID: cloudflareId, Description: "Cloudflare DNS"},
			{ID: googleCloudId, Description: "Google Cloud DNS"},
		},
	}

	awsAccessKey = dynamic_fields.DynamicField{
		ID:          "awsAccessKey",
		Priority:    2,
		Description: "AWS access key (for the DNS challenge)",
		Required:    true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       awsRoute53Id,
		},
	}

	awsSecretKey = dynamic_fields.DynamicField{
		ID:          "awsSecretKey",
		Priority:    3,
		Description: "AWS secret key",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       awsRoute53Id,
		},
	}

	awsHostedZoneID = dynamic_fields.DynamicField{
		ID:          "awsHostedZoneId",
		Priority:    3,
		Description: "AWS hosted zone ID",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       awsRoute53Id,
		},
	}

	cloudflareApiToken = dynamic_fields.DynamicField{
		ID:          "cloudflareApiToken",
		Priority:    2,
		Description: "Cloudflare API token (for the DNS challenge)",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       cloudflareId,
		},
	}

	googleCloudPrivateKey = dynamic_fields.DynamicField{
		ID:          "googleCloudPrivateKey",
		Priority:    2,
		Description: "Service account private key JSON",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.MultiLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       googleCloudId,
		},
	}

	azureTenantId = dynamic_fields.DynamicField{
		ID:          "azureTenantId",
		Priority:    2,
		Description: "Azure tenant ID (for the DNS challenge)",
		Required:    true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       azureId,
		},
	}

	azureSubscriptionId = dynamic_fields.DynamicField{
		ID:          "azureSubscriptionId",
		Priority:    3,
		Description: "Azure subscription ID",
		Required:    true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       azureId,
		},
	}

	azureClientId = dynamic_fields.DynamicField{
		ID:          "azureClientId",
		Priority:    4,
		Description: "Azure client ID",
		Required:    true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       azureId,
		},
	}

	azureClientSecret = dynamic_fields.DynamicField{
		ID:          "azureClientSecret",
		Priority:    5,
		Description: "Azure client secret",
		Required:    true,
		Sensitive:   true,
		Type:        dynamic_fields.SingleLineTextType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       azureId,
		},
	}

	azureEnvironment = dynamic_fields.DynamicField{
		ID:          "azureEnvironment",
		Priority:    6,
		Description: "Azure environment",
		Required:    true,
		Type:        dynamic_fields.EnumType,
		Condition: &dynamic_fields.Condition{
			ParentField: dnsProvider.ID,
			Value:       azureId,
		},
		EnumOptions: &[]*dynamic_fields.EnumOption{
			{ID: "DEFAULT", Description: "Azure (default)"},
			{ID: "CHINA", Description: "China"},
			{ID: "US_GOVERNMENT", Description: "US Government"},
		},
	}
)

func stringPtr(s string) *string {
	return &s
}
