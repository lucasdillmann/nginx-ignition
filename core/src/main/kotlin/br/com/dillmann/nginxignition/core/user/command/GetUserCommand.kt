package br.com.dillmann.nginxignition.core.user.command

import br.com.dillmann.nginxignition.core.user.User
import java.util.UUID

interface GetUserCommand {
    suspend fun getById(id: UUID): User?
}