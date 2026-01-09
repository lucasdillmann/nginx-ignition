package aliesa

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/aliesa"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

const (
	ramRoleFieldID       = "aliesaRamRole"
	apiKeyFieldID        = "aliesaApiKey"
	secretKeyFieldID     = "aliesaSecretKey"
	securityTokenFieldID = "aliesaSecurityToken"
	regionIDFieldID      = "aliesaRegionID"
)

type Provider struct{}

func (p *Provider) ID() string {
	return "ALIESA"
}

func (p *Provider) Name() string {
	return "AlibabaCloud ESA"
}

func (p *Provider) DynamicFields() []dynamicfields.DynamicField {
	return dns.LinkedToProvider(p.ID(), []dynamicfields.DynamicField{
		{
			ID:          ramRoleFieldID,
			Description: "AlibabaCloud RAM role",
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          apiKeyFieldID,
			Description: "AlibabaCloud RAM Access Key ID",
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          secretKeyFieldID,
			Description: "AlibabaCloud RAM Access Key Secret",
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          securityTokenFieldID,
			Description: "AlibabaCloud RAM Security Token",
			Required:    false,
			Sensitive:   true,
			Type:        dynamicfields.SingleLineTextType,
		},
		{
			ID:          regionIDFieldID,
			Description: "AlibabaCloud Region ID",
			Required:    false,
			Type:        dynamicfields.SingleLineTextType,
		},
	})
}

func (p *Provider) ChallengeProvider(
	_ context.Context,
	_ []string,
	parameters map[string]any,
) (challenge.Provider, error) {
	ramRole, _ := parameters[ramRoleFieldID].(string)
	apiKey, _ := parameters[apiKeyFieldID].(string)
	secretKey, _ := parameters[secretKeyFieldID].(string)
	securityToken, _ := parameters[securityTokenFieldID].(string)
	regionID, _ := parameters[regionIDFieldID].(string)

	cfg := aliesa.NewDefaultConfig()
	cfg.RAMRole = ramRole
	cfg.APIKey = apiKey
	cfg.SecretKey = secretKey
	cfg.SecurityToken = securityToken
	cfg.RegionID = regionID
	cfg.TTL = dns.TTL
	cfg.PropagationTimeout = dns.PropagationTimeout
	cfg.PollingInterval = dns.PollingInterval

	return aliesa.NewDNSProviderConfig(cfg)
}
