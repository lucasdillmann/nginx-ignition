//go:build windows

package configuration

import (
	"os"
	"path/filepath"
)

func init() {
	basePath := filepath.Join(os.TempDir(), "nginx-ignition")

	defaultValues["nginx-ignition.nginx.config-path"] = filepath.Join(basePath, "nginx")
	defaultValues["nginx-ignition.vpn.config-path"] = filepath.Join(basePath, "vpn")
	defaultValues["nginx-ignition.database.data-path"] = filepath.Join(basePath, "data")
}
