package br.com.dillmann.nginxignition.core.host.command

import java.util.UUID

fun interface HostExistsByIdCommand {
    suspend fun existsById(id: UUID): Boolean
}
