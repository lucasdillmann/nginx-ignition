package br.com.dillmann.nginxsidewheel.database.common

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import org.jetbrains.exposed.sql.Database

internal class DatabaseStartup: StartupCommand {
    override val priority = 0

    override fun execute() {
        Database.connect(
            url = "jdbc:h2:mem:nginx_sidewheel;DB_CLOSE_DELAY=-1",
            user = "root",
            driver = "org.h2.Driver",
            password = "",
        )
    }
}
