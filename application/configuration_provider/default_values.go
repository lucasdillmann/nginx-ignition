package configuration_provider

func defaultValues() map[string]string {
	return map[string]string{
		"nginx-ignition.server.port":                                   "8090",
		"nginx-ignition.server.shutdown-delay-seconds":                 "2",
		"nginx-ignition.nginx.binary-path":                             "nginx",
		"nginx-ignition.nginx.config-directory":                        "/tmp/nginx-ignition/nginx",
		"nginx-ignition.database.url":                                  "jdbc:h2:mem:nginx-ignition;DB_CLOSE_DELAY=-1",
		"nginx-ignition.database.username":                             "root",
		"nginx-ignition.database.password":                             "root",
		"nginx-ignition.database.connection-pool.maximum-size":         "5",
		"nginx-ignition.database.connection-pool.minimum-size":         "1",
		"nginx-ignition.security.user-password-hashing.algorithm":      "SHA-512",
		"nginx-ignition.security.user-password-hashing.salt-size":      "64",
		"nginx-ignition.security.user-password-hashing.iterations":     "1024",
		"nginx-ignition.security.jwt.secret":                           "",
		"nginx-ignition.security.jwt.ttl-seconds":                      "3600",
		"nginx-ignition.security.jwt.clock-skew-seconds":               "60",
		"nginx-ignition.security.jwt.renew-window-seconds":             "900",
		"nginx-ignition.certificate.lets-encrypt.production":           "true",
		"nginx-ignition.integration.truenas.api-cache-timeout-seconds": "15",
	}
}
