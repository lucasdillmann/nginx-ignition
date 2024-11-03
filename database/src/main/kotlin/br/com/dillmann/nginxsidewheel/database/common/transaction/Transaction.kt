package br.com.dillmann.nginxsidewheel.database.common.transaction

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.withContext
import org.jetbrains.exposed.sql.transactions.transaction

internal suspend fun <T> coTransaction(action: suspend () -> T) =
    withContext(Dispatchers.IO) {
        transaction {
            runBlocking {
                action()
            }
        }
    }
