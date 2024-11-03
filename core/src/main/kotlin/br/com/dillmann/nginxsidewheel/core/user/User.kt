package br.com.dillmann.nginxsidewheel.core.user

import java.util.UUID

data class User(
    val id: UUID,
    val enabled: Boolean,
    val name: String,
    val username: String,
    val passwordHash: String,
    val passwordSalt: String,
    val role: Role,
) {
    enum class Role {
        ADMIN,
        REGULAR_USER,
    }
}
