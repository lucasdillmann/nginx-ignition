package smtp

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/notification"
)

type Provider struct {
	sender mailSender
}

func newProvider() *Provider {
	return &Provider{sender: defaultMailSender{}}
}

func (p *Provider) ID() string {
	return providerID
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.NotificationSmtpName)
}

func (p *Provider) ImportantInstructions(ctx context.Context) []*i18n.Message {
	return importantInstructions(ctx)
}

func (p *Provider) ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return configurationFields(ctx)
}

func (p *Provider) Send(
	ctx context.Context,
	parameters map[string]any,
	deliverable notification.Deliverable,
) error {
	settings, recipients := parseMailSettings(parameters)

	message := buildMessage(deliverable, settings.fromAddress, recipients)

	return p.sender.Send(ctx, settings, recipients, message)
}
