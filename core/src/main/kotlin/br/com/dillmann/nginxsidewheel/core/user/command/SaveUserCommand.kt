package br.com.dillmann.nginxsidewheel.core.user.command

import br.com.dillmann.nginxsidewheel.core.user.model.SaveUserRequest

interface SaveUserCommand {
    suspend fun save(request: SaveUserRequest)
}
