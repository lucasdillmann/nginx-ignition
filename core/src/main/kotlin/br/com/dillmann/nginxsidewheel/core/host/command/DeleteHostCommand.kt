package br.com.dillmann.nginxsidewheel.core.host.command

import java.util.UUID

interface DeleteHostCommand {
    suspend fun deleteById(id: UUID)
}
