package constants

import "regexp"

var TLDPattern = regexp.MustCompile(`^(?:[a-zA-Z0-9*](?:[a-zA-Z0-9-*]{0,61}[a-zA-Z0-9*])?\.)+[a-zA-Z]{2,}$`)
