package br.com.dillmann.nginxignition.core.host.command

import br.com.dillmann.nginxignition.core.host.Host
import java.util.UUID

fun interface GetHostCommand {
    suspend fun getById(id: UUID): Host?
}
