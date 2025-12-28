package resolver

import (
	"regexp"
)

const (
	defaultDockerDNSIP = "127.0.0.11"
	hostQualifier      = "host"
	containerQualifier = "container"
	ingressQualifier   = "ingress"
	httpURLTemplate    = "http://%s:%d"
)

var (
	containerNameGeneralNormalizationRegex    = regexp.MustCompile(`[^a-zA-Z0-9\-_.]+`)
	containerNameUnderscoreNormalizationRegex = regexp.MustCompile(`_{2,}`)
)
