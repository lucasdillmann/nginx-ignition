package br.com.dillmann.nginxignition.application.configuration

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider

internal class RootConfigurationProvider: ConfigurationProvider {
    override fun get(key: String): String =
        resolveFromEnvironmentVariables(key) ?: ConfigurationDefaults[key] ?: error("Key not found: $key")

    override fun withPrefix(prefix: String): ConfigurationProvider =
        PrefixedConfigurationProvider(this, prefix)

    private fun resolveFromEnvironmentVariables(key: String): String? =
        System.getenv(key) ?: System.getenv(key.toSnakeCase())

    private fun String.toSnakeCase() =
        replace("-", "_").replace(".", "_").uppercase()
}
