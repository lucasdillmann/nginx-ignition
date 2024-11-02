package br.com.dillmann.nginxsidewheel.database.host

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostTable
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction

internal class HostMigrations: StartupCommand {
    override val priority: Int = 100

    override suspend fun execute() {
        transaction {
            SchemaUtils.create(
                HostTable,
                HostBindingTable,
                HostRouteTable,
            )
        }
    }
}
