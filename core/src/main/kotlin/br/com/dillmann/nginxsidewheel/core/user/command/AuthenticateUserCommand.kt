package br.com.dillmann.nginxsidewheel.core.user.command

import br.com.dillmann.nginxsidewheel.core.user.User

interface AuthenticateUserCommand {
    suspend fun authenticate(username: String, password: String): User?
}
