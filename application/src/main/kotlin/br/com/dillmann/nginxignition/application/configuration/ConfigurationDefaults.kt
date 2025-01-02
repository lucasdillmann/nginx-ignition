package br.com.dillmann.nginxignition.application.configuration

internal val ConfigurationDefaults = mapOf(
    "nginx-ignition.server.port" to "8090",
    "nginx-ignition.server.shutdown-delay-seconds" to "2",
    "nginx-ignition.nginx.binary-path" to "nginx",
    "nginx-ignition.nginx.config-directory" to "/tmp/nginx-ignition/nginx",
    "nginx-ignition.database.url" to "jdbc:h2:mem:nginx-ignition;DB_CLOSE_DELAY=-1",
    "nginx-ignition.database.username" to "root",
    "nginx-ignition.database.password" to "root",
    "nginx-ignition.database.connection-pool.maximum-size" to "5",
    "nginx-ignition.database.connection-pool.minimum-size" to "1",
    "nginx-ignition.security.user-password-hashing.algorithm" to "SHA-512",
    "nginx-ignition.security.user-password-hashing.salt-size" to "64",
    "nginx-ignition.security.user-password-hashing.iterations" to "1024",
    "nginx-ignition.security.jwt.secret" to "",
    "nginx-ignition.security.jwt.ttl-seconds" to "3600",
    "nginx-ignition.security.jwt.clock-skew-seconds" to "60",
    "nginx-ignition.security.jwt.renew-window-seconds" to "900",
    "nginx-ignition.certificate.lets-encrypt.production" to "true",
    "nginx-ignition.integration.truenas.api-cache-timeout-seconds" to "15",
)
