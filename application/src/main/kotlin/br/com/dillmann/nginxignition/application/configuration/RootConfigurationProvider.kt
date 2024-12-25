package br.com.dillmann.nginxignition.application.configuration

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import org.yaml.snakeyaml.Yaml

class RootConfigurationProvider: ConfigurationProvider {
    private val values: Map<String, Any?>

    init {
        val stream = javaClass.classLoader.getResourceAsStream("application.yaml")
        values = Yaml().loadAs(stream, Map::class.java)
    }

    override fun get(key: String): String =
        resolveFromEnvironmentVariables(key) ?: resolveFromConfigurationFile(key)

    override fun withPrefix(prefix: String): ConfigurationProvider =
        PrefixedConfigurationProvider(this, prefix)

    private fun resolveFromConfigurationFile(key: String): String =
        try {
            getValue(key.split(".").iterator(), values)
        } catch (ex: Exception) {
            throw IllegalArgumentException("Unable to resolve configuration value for key [$key]", ex)
        }

    private tailrec fun getValue(keys: Iterator<String>, node: Any): String =
        if (keys.hasNext()) {
            val child = (node as Map<*, *>)[keys.next()]
            getValue(keys, child!!)
        } else {
            node.toString()
        }

    private fun resolveFromEnvironmentVariables(key: String): String? =
        System.getenv(key) ?: System.getenv(key.toSnakeCase())

    private fun String.toSnakeCase() =
        replace("-", "_").replace(".", "_").uppercase()
}
