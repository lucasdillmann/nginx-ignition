package br.com.dillmann.nginxignition.core.common.configuration

interface ConfigurationProvider {
    fun get(key: String): String
    fun withPrefix(prefix: String): ConfigurationProvider
}
