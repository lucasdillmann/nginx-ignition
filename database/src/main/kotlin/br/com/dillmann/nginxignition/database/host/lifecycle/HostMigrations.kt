package br.com.dillmann.nginxignition.database.host.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxignition.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxignition.database.host.mapping.HostTable
import org.jetbrains.exposed.sql.SchemaUtils

internal class HostMigrations: StartupCommand {
    override val priority: Int = 100

    override suspend fun execute() {
        coTransaction {
            SchemaUtils.createMissingTablesAndColumns(
                HostTable,
                HostBindingTable,
                HostRouteTable,
            )
        }
    }
}
