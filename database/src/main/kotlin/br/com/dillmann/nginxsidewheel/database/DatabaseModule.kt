package br.com.dillmann.nginxsidewheel.database

import br.com.dillmann.nginxsidewheel.core.host.HostRepository
import br.com.dillmann.nginxsidewheel.database.host.HostDatabaseRepository
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostTable
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction
import org.koin.core.module.Module
import org.koin.dsl.module

object DatabaseModule {
    fun initialize(): Module {
        // TODO: Replace with some more final solution
        Database.connect(
            url = "jdbc:h2:mem:nginx_sidewheel;DB_CLOSE_DELAY=-1",
            user = "root",
            driver = "org.h2.Driver",
            password = "",
        )

        transaction {
            SchemaUtils.create(
                HostTable,
                HostBindingTable,
                HostRouteTable,
            )
        }

        return module {
            single<HostRepository> { HostDatabaseRepository() }
        }
    }
}
