package br.com.dillmann.nginxsidewheel.database.common

import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import org.jetbrains.exposed.sql.Database

internal class DatabaseStartup(private val configurationProvider: ConfigurationProvider): StartupCommand {
    override val priority = 0

    override suspend fun execute() {
        val url = configurationProvider.get("nginx-sidewheel.database.url")
        val username = configurationProvider.get("nginx-sidewheel.database.username")
        val password = configurationProvider.get("nginx-sidewheel.database.password")

        logger<DatabaseStartup>().info("Starting database connection to $url with username $username")

        Database.connect(
            url = url,
            user = username,
            password = password,
        )
    }
}
