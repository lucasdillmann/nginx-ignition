package notification

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type testProvider struct{}

func (testProvider) ID() string { return "SMTP" }

func (testProvider) Name(_ context.Context) *i18n.Message {
	return i18n.Static("Test")
}

func (testProvider) ImportantInstructions(_ context.Context) []*i18n.Message {
	return nil
}

func (testProvider) ConfigurationFields(_ context.Context) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:        "host",
			Sensitive: false,
			Type:      dynamicfields.SingleLineTextType,
		},
		{
			ID:        "password",
			Sensitive: true,
			Type:      dynamicfields.SingleLineTextType,
		},
	}
}

func (testProvider) Send(
	_ context.Context,
	_ map[string]any,
	_ Deliverable,
) error {
	return nil
}

func newConfiguration() *Configuration {
	return &Configuration{
		ID:         uuid.New(),
		Name:       "test",
		Provider:   "SMTP",
		Enabled:    true,
		Parameters: map[string]any{},
	}
}

func testProviders() []Provider {
	return []Provider{testProvider{}}
}

type failingTestProvider struct {
	testProvider
	sendError error
}

func (provider failingTestProvider) Send(
	_ context.Context,
	_ map[string]any,
	_ Deliverable,
) error {
	return provider.sendError
}

func failingTestProviders(sendError error) func() []Provider {
	return func() []Provider {
		return []Provider{failingTestProvider{sendError: sendError}}
	}
}

type requiredHostProvider struct {
	testProvider
}

func (requiredHostProvider) ConfigurationFields(
	_ context.Context,
) []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		{
			ID:       "host",
			Required: true,
			Type:     dynamicfields.SingleLineTextType,
		},
	}
}

var errSendFailed = errors.New("send failed")
