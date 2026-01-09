//nolint:misspell
package internetbs

//nolint:misspell
import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/internetbs"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

//nolint:gosec
const (
	apiKeyFieldID   = "internetBsApiKey"
	passwordFieldID = "internetBsPassword"
)

type Provider struct{}

//nolint:misspell
func (p *Provider) ID() string { return "INTERNETBS" }

func (p *Provider) Name() string { return "Internet.bs" }

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          apiKeyFieldID,
			Description: "Internet.bs API key",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          passwordFieldID,
			Description: "Internet.bs password",
			Required:    true,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	apiKey, _ := parameters[apiKeyFieldID].(string)
	password, _ := parameters[passwordFieldID].(string)

	//nolint:misspell
	cfg := internetbs.NewDefaultConfig()
	cfg.APIKey = apiKey
	cfg.Password = password
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval
	cfg.TTL = dns.TTL

	//nolint:misspell
	return internetbs.NewDNSProviderConfig(cfg)
}
