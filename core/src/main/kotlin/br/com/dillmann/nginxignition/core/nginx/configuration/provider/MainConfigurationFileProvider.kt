package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider

internal class MainConfigurationFileProvider: NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> {
        val contents = """
            worker_processes 2;
            error_log $basePath/logs/main.log;
            pid $basePath/nginx.pid;
            
            events {
                worker_connections 1024;
            }
            
            http {
                default_type application/octet-stream;
                sendfile on;
                keepalive_timeout 65;
                server_tokens off;
                client_max_body_size 1024G;
                
                include $basePath/config/mime.types;
                ${hosts.joinToString(separator = "\n") { 
                    "include $basePath/config/host-${it.id}.conf;"
                }}
            }
        """.trimIndent()

        return listOf(
            NginxConfigurationFileProvider.Output(
                name = "nginx.conf",
                contents = contents,
            )
        )
    }
}
