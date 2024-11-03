package br.com.dillmann.nginxignition.core.nginx.exception

class NginxCommandException(
    val command: String,
    val exitCode: Int,
    val output: String,
) : RuntimeException(
    "Nginx command failed with exit code $exitCode: $output"
)
