package br.com.dillmann.nginxignition.core.user.command

import br.com.dillmann.nginxignition.core.user.model.SaveUserRequest

interface SaveUserCommand {
    suspend fun save(request: SaveUserRequest)
}
