package br.com.dillmann.nginxsidewheel.core.user.command

import java.util.UUID

interface DeleteUserCommand {
    suspend fun deleteById(id: UUID)
}
