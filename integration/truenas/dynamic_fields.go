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
			Description: i18n.M(ctx, i18n.K.TruenasCommonUrl),
			Priority:    1,
			Required:    true,
			HelpText:    i18n.M(ctx, i18n.K.TruenasCommonUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:          proxyURLFieldID,
			Description: i18n.M(ctx, i18n.K.TruenasCommonProxyUrl),
			Priority:    2,
			Required:    false,
			HelpText:    i18n.M(ctx, i18n.K.TruenasCommonProxyUrlHelp),
			Type:        dynamicfields.URLType,
		},
		{
			ID:          usernameFieldID,
			Description: i18n.M(ctx, i18n.K.TruenasCommonUsername),
			Priority:    3,
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: i18n.M(ctx, i18n.K.TruenasCommonPassword),
			Priority:    4,
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	}
}
