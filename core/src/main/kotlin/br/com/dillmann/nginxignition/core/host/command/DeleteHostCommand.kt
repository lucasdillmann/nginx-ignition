package br.com.dillmann.nginxignition.core.host.command

import java.util.UUID

fun interface DeleteHostCommand {
    suspend fun deleteById(id: UUID)
}
