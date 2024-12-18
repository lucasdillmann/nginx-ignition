package br.com.dillmann.nginxignition.core.accesslist.command

import java.util.*

fun interface DeleteAccessListByIdCommand {
    suspend fun deleteById(id: UUID)
}
