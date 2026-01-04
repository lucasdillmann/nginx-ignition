package nginx

import "dillmann.com.br/nginx-ignition/core/nginx"

func newMetadata() *nginx.Metadata {
	return &nginx.Metadata{
		Version:      "1.21.0",
		BuildDetails: "nginx version: nginx/1.21.0",
		Modules:      []string{"http_ssl_module", "http_v2_module"},
	}
}
