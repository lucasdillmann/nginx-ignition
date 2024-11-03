package br.com.dillmann.nginxignition.core.host.command

import java.util.UUID

interface DeleteHostCommand {
    suspend fun deleteById(id: UUID)
}
