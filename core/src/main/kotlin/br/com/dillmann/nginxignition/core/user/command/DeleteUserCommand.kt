package br.com.dillmann.nginxignition.core.user.command

import java.util.UUID

fun interface DeleteUserCommand {
    suspend fun deleteById(id: UUID)
}
