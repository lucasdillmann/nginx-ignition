package br.com.dillmann.nginxignition.api.user.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.user.User.Role
import kotlinx.serialization.Serializable

@Serializable
internal data class UserResponse(
    val id: UuidString,
    val enabled: Boolean,
    val name: String,
    val username: String,
    val role: Role,
)
