package br.com.dillmann.nginxignition.core.user.command

import java.util.UUID

interface UpdateUserPasswordCommand {
    suspend fun updatePassword(userId: UUID, currentPassword: String, newPassword: String)
}
