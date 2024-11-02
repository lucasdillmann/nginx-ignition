package br.com.dillmann.nginxsidewheel.application.common.provider

import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import io.ktor.server.config.*

class CompositeConfigurationProvider: ConfigurationProvider {
    private val ktorDelegate = ConfigLoader.load()

    override fun get(key: String): String =
        resolveFromEnvironmentVariables(key) ?: ktorDelegate.property(key).getString()

    private fun resolveFromEnvironmentVariables(key: String): String? =
        System.getenv(key) ?: System.getenv(key.toSnakeCase())

    private fun String.toSnakeCase() =
        replace("-", "_").replace(".", "_").uppercase()
}
