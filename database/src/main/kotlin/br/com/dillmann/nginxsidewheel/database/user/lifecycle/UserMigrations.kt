package br.com.dillmann.nginxsidewheel.database.user.lifecycle

import br.com.dillmann.nginxsidewheel.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxsidewheel.database.common.transaction.coTransaction
import br.com.dillmann.nginxsidewheel.database.user.mapping.UserTable
import org.jetbrains.exposed.sql.SchemaUtils

class UserMigrations: StartupCommand {
    override val priority = 100

    override suspend fun execute() {
        coTransaction {
            SchemaUtils.createMissingTablesAndColumns(UserTable)
        }
    }
}
