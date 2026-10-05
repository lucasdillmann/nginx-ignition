package constants

import "regexp"

var TLDPattern = regexp.MustCompile(
	`^(?:[a-zA-Z0-9*](?:[a-zA-Z0-9-*]{0,61}[a-zA-Z0-9*])?\.)+[a-zA-Z]{2,}$`,
)

var HostnamePattern = regexp.MustCompile(
	`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
)

var SubdomainNamePattern = regexp.MustCompile(
	`^[a-z](?:[a-z0-9_-]{0,61}[a-z0-9])?$`,
)
