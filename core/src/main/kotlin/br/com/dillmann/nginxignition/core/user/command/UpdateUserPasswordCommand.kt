package br.com.dillmann.nginxignition.core.user.command

import java.util.UUID

fun interface UpdateUserPasswordCommand {
    suspend fun updatePassword(userId: UUID, currentPassword: String, newPassword: String)
}
