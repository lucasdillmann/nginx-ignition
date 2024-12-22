package br.com.dillmann.nginxignition.core.accesslist.command

import java.util.*

fun interface DeleteAccessListByIdCommand {
    data class Output(
        val deleted: Boolean,
        val reason: String,
    )

    suspend fun deleteById(id: UUID): Output
}
