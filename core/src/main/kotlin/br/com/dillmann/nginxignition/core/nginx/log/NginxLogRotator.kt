package br.com.dillmann.nginxignition.core.nginx.log

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.host.HostRepository
import br.com.dillmann.nginxignition.core.nginx.NginxProcessManager
import br.com.dillmann.nginxignition.core.settings.SettingsRepository

@Suppress("UnusedPrivateProperty") // TODO: Remove this
internal class NginxLogRotator(
    private val configurationProvider: ConfigurationProvider,
    private val settingsRepository: SettingsRepository,
    private val hostRepository: HostRepository,
    private val processManager: NginxProcessManager,
) {
    fun rotate() {
        // TODO: Implement this
    }
}
