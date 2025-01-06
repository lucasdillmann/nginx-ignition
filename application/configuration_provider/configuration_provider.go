package configuration_provider

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
)

type provider struct {
	prefix string
}

func New() configuration.Configuration {
	return &provider{}
}

func (p *provider) Get(_ string) (string, error) {
	return "", core_errors.NotImplemented()
}

func (p *provider) WithPrefix(prefix string) configuration.Configuration {
	var newPrefix string
	if p.prefix == "" {
		newPrefix = prefix
	} else {
		newPrefix = p.prefix + "." + prefix
	}

	return &provider{newPrefix}
}
