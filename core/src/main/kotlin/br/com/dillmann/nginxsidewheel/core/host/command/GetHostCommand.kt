package br.com.dillmann.nginxsidewheel.core.host.command

import br.com.dillmann.nginxsidewheel.core.host.Host
import java.util.UUID

interface GetHostCommand {
    suspend fun getById(id: UUID): Host?
}
