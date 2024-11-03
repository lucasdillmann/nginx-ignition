package br.com.dillmann.nginxignition.core.user.model

import br.com.dillmann.nginxignition.core.user.User.Role
import java.util.*

data class SaveUserRequest(
    val id: UUID,
    val enabled: Boolean,
    val name: String,
    val username: String,
    val password: String?,
    val role: Role,
)
