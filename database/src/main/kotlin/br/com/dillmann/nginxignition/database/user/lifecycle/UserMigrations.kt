package br.com.dillmann.nginxignition.database.user.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.user.mapping.UserTable
import org.jetbrains.exposed.sql.SchemaUtils

class UserMigrations: StartupCommand {
    override val priority = 100

    override suspend fun execute() {
        coTransaction {
            SchemaUtils.createMissingTablesAndColumns(UserTable)
        }
    }
}
