package br.com.dillmann.nginxsidewheel.core.common.provider

interface ConfigurationProvider {
    fun get(key: String): String
    fun withPrefix(prefix: String): ConfigurationProvider
}
