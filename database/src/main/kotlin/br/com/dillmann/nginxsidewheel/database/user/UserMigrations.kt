package br.com.dillmann.nginxsidewheel.database.user

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.database.common.coTransaction
import br.com.dillmann.nginxsidewheel.database.user.mapping.UserTable
import org.jetbrains.exposed.sql.SchemaUtils

class UserMigrations: StartupCommand {
    override val priority = 100

    override suspend fun execute() {
        coTransaction {
            SchemaUtils.create(UserTable)
        }
    }
}
