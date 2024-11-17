package br.com.dillmann.nginxignition.database.common.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider
import com.zaxxer.hikari.HikariDataSource
import org.h2.Driver as H2Driver
import org.jetbrains.exposed.sql.Database
import org.postgresql.Driver as PostgreSqlDriver
import org.mariadb.jdbc.Driver as MariaDbDriver

internal class DatabaseStartup(configurationProvider: ConfigurationProvider): StartupCommand {
    private val logger = logger<DatabaseStartup>()
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.database")
    override val priority = 0

    override suspend fun execute() {
        val url = configurationProvider.get("url")
        val username = configurationProvider.get("username")

        logger.info("Starting database connection to $url with username $username")

        val dataSource = HikariDataSource()
        dataSource.jdbcUrl = url
        dataSource.username = username
        dataSource.password = configurationProvider.get("password")
        dataSource.maximumPoolSize = configurationProvider.get("connection-pool.maximum-size").toInt()
        dataSource.minimumIdle = configurationProvider.get("connection-pool.minimum-size").toInt()
        dataSource.driverClassName = resolveDriverClassName(url)

        Database.connect(dataSource)
    }

    private fun resolveDriverClassName(url: String): String? =
        when {
            url.startsWith("jdbc:mysql:") || url.startsWith("jdbc:mariadb:") ->
                MariaDbDriver::class.qualifiedName

            url.startsWith("jdbc:postgresql:") ->
                PostgreSqlDriver::class.qualifiedName

            url.startsWith("jdbc:h2:") -> {
                logger.warn(
                    "Application is configured to use the embedded H2 database. This isn't recommended for " +
                        "production environments, please refer to the documentation in order to migrate to another " +
                        "DBMS like PostgreSQL or MariaDB/MySQL."
                )
                H2Driver::class.qualifiedName
            }

            else -> null
        }
}
