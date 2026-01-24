package truenas

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	urlFieldID      = "url"
	proxyURLFieldID = "proxyUrl"
	usernameFieldID = "username"
	passwordFieldID = "password"
)

func dynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:          urlFieldID,
			Description: i18n.M(ctx, i18n.K.IntegrationTruenasUrl),
			Priority:    1,
			Required:    true,
			HelpText:    i18n.M(ctx, i18n.K.IntegrationTruenasUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:          proxyURLFieldID,
			Description: i18n.M(ctx, i18n.K.IntegrationTruenasProxyUrl),
			Priority:    2,
			Required:    false,
			HelpText:    i18n.M(ctx, i18n.K.IntegrationTruenasProxyUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.CommonUsername),
			Priority:    3,
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.IntegrationTruenasPassword),
			Priority:    4,
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	}
}
