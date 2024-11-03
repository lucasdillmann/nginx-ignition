package br.com.dillmann.nginxsidewheel.core.user.command

import br.com.dillmann.nginxsidewheel.core.user.User
import java.util.UUID

interface GetUserCommand {
    suspend fun getById(id: UUID): User?
}
