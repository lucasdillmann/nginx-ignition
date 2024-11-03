package br.com.dillmann.nginxignition.core.host.command

import br.com.dillmann.nginxignition.core.host.Host

interface SaveHostCommand {
    suspend fun save(input: Host)

}
