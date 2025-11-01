package vpn

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
)

var (
	ErrVpnNotFound = core_error.New("VPN not found", true)
	ErrVpnDisabled = core_error.New("VPN is disabled", true)
)
