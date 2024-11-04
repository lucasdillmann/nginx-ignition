package br.com.dillmann.nginxignition.core.user.command

import java.util.UUID

interface GetUserStatusCommand {
    suspend fun isEnabled(id: UUID): Boolean
}
