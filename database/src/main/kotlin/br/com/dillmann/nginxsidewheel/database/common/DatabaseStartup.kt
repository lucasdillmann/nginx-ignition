package br.com.dillmann.nginxsidewheel.database.common

import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import org.jetbrains.exposed.sql.Database

internal class DatabaseStartup(configurationProvider: ConfigurationProvider): StartupCommand {
    private val configurationProvider = configurationProvider.withPrefix("nginx-sidewheel.database")
    override val priority = 0

    override suspend fun execute() {
        val url = configurationProvider.get("url")
        val username = configurationProvider.get("username")
        val password = configurationProvider.get("password")

        logger<DatabaseStartup>().info("Starting database connection to $url with username $username")

        Database.connect(
            url = url,
            user = username,
            password = password,
        )
    }
}
