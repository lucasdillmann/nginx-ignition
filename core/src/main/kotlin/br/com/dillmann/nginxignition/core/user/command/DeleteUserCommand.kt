package br.com.dillmann.nginxignition.core.user.command

import java.util.UUID

interface DeleteUserCommand {
    suspend fun deleteById(id: UUID)
}
