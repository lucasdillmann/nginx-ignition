package br.com.dillmann.nginxignition.core.host.command

import br.com.dillmann.nginxignition.core.host.Host

fun interface SaveHostCommand {
    suspend fun save(input: Host)

}
