package configuration

func defaultValues() map[string]string {
	return map[string]string{
		"nginx-ignition.server.binding-port":                           "8090",
		"nginx-ignition.server.binding-address":                        "0.0.0.0",
		"nginx-ignition.nginx.binary-path":                             "nginx",
		"nginx-ignition.nginx.config-path":                             "/tmp/nginx-ignition/nginx",
		"nginx-ignition.database.driver":                               "sqlite",
		"nginx-ignition.database.data-path":                            "/tmp/nginx-ignition/data",
		"nginx-ignition.security.user-password-hashing.algorithm":      "SHA-512",
		"nginx-ignition.security.user-password-hashing.salt-size":      "64",
		"nginx-ignition.security.user-password-hashing.iterations":     "1024",
		"nginx-ignition.security.jwt.secret":                           "",
		"nginx-ignition.security.jwt.ttl-seconds":                      "3600",
		"nginx-ignition.security.jwt.clock-skew-seconds":               "60",
		"nginx-ignition.security.jwt.renew-window-seconds":             "900",
		"nginx-ignition.certificate.lets-encrypt.production":           "true",
		"nginx-ignition.integration.truenas.api-cache-timeout-seconds": "15",
		"nginx-ignition.password-reset.username":                       "",
	}
}
