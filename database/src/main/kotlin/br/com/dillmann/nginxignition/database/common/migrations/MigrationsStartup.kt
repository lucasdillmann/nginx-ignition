package br.com.dillmann.nginxignition.database.common.migrations

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand

internal class MigrationsStartup(private val service: MigrationsService): StartupCommand {
    override suspend fun execute() {
        service.migrateDatabase()
    }
}
