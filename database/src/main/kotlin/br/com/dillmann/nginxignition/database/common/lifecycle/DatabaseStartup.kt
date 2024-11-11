package br.com.dillmann.nginxignition.database.common.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider
import com.zaxxer.hikari.HikariDataSource
import org.jetbrains.exposed.sql.Database

internal class DatabaseStartup(configurationProvider: ConfigurationProvider): StartupCommand {
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.database")
    override val priority = 0

    override suspend fun execute() {
        val url = configurationProvider.get("url")
        val username = configurationProvider.get("username")

        logger<DatabaseStartup>().info("Starting database connection to $url with username $username")

        val dataSource = HikariDataSource()
        dataSource.jdbcUrl = url
        dataSource.username = username
        dataSource.password = configurationProvider.get("password")
        dataSource.maximumPoolSize = configurationProvider.get("connection-pool.maximum-size").toInt()
        dataSource.minimumIdle = configurationProvider.get("connection-pool.minimum-size").toInt()

        Database.connect(dataSource)
    }
}
