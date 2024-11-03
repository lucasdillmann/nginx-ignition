package br.com.dillmann.nginxignition.application.common.provider

import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider

internal class PrefixedConfigurationProvider(
    private val delegate: ConfigurationProvider,
    private val prefix: String,
): ConfigurationProvider {
    override fun get(key: String): String =
        delegate.get("$prefix.$key")

    override fun withPrefix(prefix: String): ConfigurationProvider =
        PrefixedConfigurationProvider(this, prefix)
}
