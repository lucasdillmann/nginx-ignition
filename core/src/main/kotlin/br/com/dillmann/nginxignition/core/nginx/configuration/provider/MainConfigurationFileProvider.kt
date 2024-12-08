package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider
import br.com.dillmann.nginxignition.core.settings.SettingsRepository

internal class MainConfigurationFileProvider(
    private val settingsService: SettingsRepository,
): NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> {
        val settings = settingsService.get().nginx
        val logs = settings.logs
        val contents = """
            worker_processes ${settings.workerProcesses};
            pid $basePath/nginx.pid;
            error_log ${
                if (logs.serverLogsEnabled) "$basePath/logs/main.log ${logs.serverLogsLevel.name.lowercase()}" 
                else "off"
            };
            
            events {
                worker_connections ${settings.workerConnections};
            }
            
            http {
                default_type ${settings.defaultContentType};
                sendfile ${enabledFlag(settings.sendfileEnabled)};
                server_tokens ${enabledFlag(settings.serverTokensEnabled)};
                client_max_body_size ${settings.maximumBodySizeMb}M;
                
                keepalive_timeout ${settings.timeouts.keepalive};
                proxy_connect_timeout ${settings.timeouts.connect};
                proxy_read_timeout ${settings.timeouts.read};
                proxy_send_timeout ${settings.timeouts.send};
                send_timeout ${settings.timeouts.send};
                
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

    private fun enabledFlag(value: Boolean) = if (value) "on" else "off"
}
