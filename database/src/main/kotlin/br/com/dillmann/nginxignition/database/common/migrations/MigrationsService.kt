package br.com.dillmann.nginxignition.database.common.migrations

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.database.common.database.DatabaseState
import br.com.dillmann.nginxignition.database.common.database.DatabaseType
import org.flywaydb.core.Flyway
import org.flywaydb.core.api.output.MigrateResult

internal class MigrationsService {
    private val logger = logger<MigrationsService>()

    fun migrateDatabase() {
        val result = runMigrations()
        with (result) {
            val targetVersion = targetSchemaVersion
            val sourceVersion = initialSchemaVersion

            when {
                !success ->
                    error("Database upgrade failed: $exception")
                successfulMigrations.isEmpty() ->
                    logger.info("Database is already up-to-date. No upgrades needed.")
                sourceVersion == null ->
                    logger.info("Database successfully initialized to version $targetVersion")
                else ->
                    logger.info("Database upgraded successfully from version $sourceVersion to $targetVersion")

            }
        }
    }

    private fun runMigrations(): MigrateResult {
        val databaseId = resolveDatabaseId()
        return Flyway
            .configure()
            .dataSource(DatabaseState.dataSource)
            .createSchemas(true)
            .cleanDisabled(true)
            .validateOnMigrate(true)
            .installedBy("nginx-ignition")
            .table("schema_version")
            .sqlMigrationSeparator("_")
            .sqlMigrationPrefix("v")
            .failOnMissingLocations(true)
            .locations("classpath:migrations/$databaseId")
            .executeInTransaction(true)
            .communityDBSupportEnabled(true)
            .load()
            .migrate()
    }

    private fun resolveDatabaseId() =
        when (DatabaseState.type) {
            DatabaseType.POSTGRESQL -> "postgresql"
            DatabaseType.H2 -> "h2"
        }
}
