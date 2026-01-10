//go:build windows

package configuration

import (
	"path/filepath"
)

func init() {
	defaultValues["nginx-ignition.nginx.config-path"] =
		filepath.Join("C:", "ProgramData", "nginx-ignition", "nginx")
	defaultValues["nginx-ignition.vpn.config-path"] =
		filepath.Join("C:", "ProgramData", "nginx-ignition", "vpn")
	defaultValues["nginx-ignition.database.data-path"] =
		filepath.Join("C:", "ProgramData", "nginx-ignition", "data")
}
