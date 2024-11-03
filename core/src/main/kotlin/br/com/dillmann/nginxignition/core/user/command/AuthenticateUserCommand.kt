package br.com.dillmann.nginxignition.core.user.command

import br.com.dillmann.nginxignition.core.user.User

interface AuthenticateUserCommand {
    suspend fun authenticate(username: String, password: String): User?
}