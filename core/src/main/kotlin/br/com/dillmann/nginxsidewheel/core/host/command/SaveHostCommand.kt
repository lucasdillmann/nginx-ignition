package br.com.dillmann.nginxsidewheel.core.host.command

import br.com.dillmann.nginxsidewheel.core.host.Host

interface SaveHostCommand {
    suspend fun save(input: Host)

}
