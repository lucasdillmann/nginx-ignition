package br.com.dillmann.nginxignition.database.common.database

import org.jetbrains.exposed.sql.Database
import javax.sql.DataSource

internal object DatabaseState {
    lateinit var database: Database
        private set

    lateinit var dataSource: DataSource
        private set

    lateinit var type: DatabaseType
        private set

    fun init(dataSource: DataSource, database: Database, type: DatabaseType) {
        if (this::database.isInitialized)
            error("Internal state was already initialized")

        DatabaseState.database = database
        DatabaseState.dataSource = dataSource
        DatabaseState.type = type
    }
}
