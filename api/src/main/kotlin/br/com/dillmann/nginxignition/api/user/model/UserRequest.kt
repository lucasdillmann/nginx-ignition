package br.com.dillmann.nginxignition.api.user.model

import br.com.dillmann.nginxignition.core.user.User.Role
import kotlinx.serialization.Serializable

@Serializable
internal data class UserRequest(
    val enabled: Boolean,
    val name: String,
    val username: String,
    val password: String? = null,
    val role: Role,
)
