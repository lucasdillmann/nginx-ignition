package br.com.dillmann.nginxignition.database.common.lifecycle

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import org.jetbrains.exposed.sql.Database

internal class DatabaseStartup(configurationProvider: ConfigurationProvider): StartupCommand {
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.database")
    override val priority = 0

    override suspend fun execute() {
        val url = configurationProvider.get("url")
        val username = configurationProvider.get("username")
        val password = configurationProvider.get("password")

        logger<DatabaseStartup>().info("Starting database connection to $url with username $username")

        // TODO: Migrate to a pooled connection source
        Database.connect(
            url = url,
            user = username,
            password = password,
        )
    }
}
