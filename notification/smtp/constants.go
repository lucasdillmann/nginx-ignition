package smtp

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const providerID = "SMTP"

const (
	hostFieldID        = "host"
	portFieldID        = "port"
	usernameFieldID    = "username"
	passwordFieldID    = "password"
	fromFieldID        = "from"
	toFieldID          = "to"
	useTLSFieldID      = "useTls"
	useStartTLSFieldID = "useStartTls"
)

func configurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:          hostFieldID,
			Priority:    0,
			Description: i18n.M(ctx, i18n.K.NotificationSmtpHost),
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:           portFieldID,
			Priority:     1,
			Description:  i18n.M(ctx, i18n.K.NotificationSmtpPort),
			Required:     true,
			Sensitive:    false,
			Type:         dynamicfields.SingleLineTextType,
			DefaultValue: "587",
		},
		{
			ID:          usernameFieldID,
			Priority:    2,
			Description: i18n.M(ctx, i18n.K.NotificationSmtpUsername),
			Required:    false,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Priority:    3,
			Description: i18n.M(ctx, i18n.K.NotificationSmtpPassword),
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          fromFieldID,
			Priority:    4,
			Description: i18n.M(ctx, i18n.K.NotificationSmtpFrom),
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.EmailType,
		},
		{
			ID:          toFieldID,
			Priority:    5,
			Description: i18n.M(ctx, i18n.K.NotificationSmtpTo),
			Required:    true,
			Sensitive:   false,
			Type:        dynamicfields.SingleLineTextType,
			HelpText:    i18n.M(ctx, i18n.K.NotificationSmtpToHelp),
		},
		{
			ID:           useTLSFieldID,
			Priority:     6,
			Description:  i18n.M(ctx, i18n.K.NotificationSmtpUseTls),
			Required:     false,
			Sensitive:    false,
			Type:         dynamicfields.BooleanType,
			DefaultValue: false,
		},
		{
			ID:           useStartTLSFieldID,
			Priority:     7,
			Description:  i18n.M(ctx, i18n.K.NotificationSmtpUseStartTls),
			Required:     false,
			Sensitive:    false,
			Type:         dynamicfields.BooleanType,
			DefaultValue: true,
			Conditions: []dynamicfields.Condition{
				{ParentField: useTLSFieldID, Value: false},
			},
		},
	}
}

func importantInstructions(ctx context.Context) []*i18n.Message {
	return []*i18n.Message{
		i18n.M(ctx, i18n.K.NotificationSmtpInstructionAppPassword),
		i18n.M(ctx, i18n.K.NotificationSmtpInstructionTls),
	}
}
