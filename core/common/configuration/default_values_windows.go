//go:build windows

package configuration

import (
	"os"
	"path/filepath"
)

func init() {
	defaultValues["nginx-ignition.nginx.config-path"] =
		filepath.Join(os.TempDir(), "nginx-ignition", "nginx")
	defaultValues["nginx-ignition.vpn.config-path"] =
		filepath.Join(os.TempDir(), "nginx-ignition", "vpn")
	defaultValues["nginx-ignition.database.data-path"] =
		filepath.Join(os.TempDir(), "nginx-ignition", "data")
}
