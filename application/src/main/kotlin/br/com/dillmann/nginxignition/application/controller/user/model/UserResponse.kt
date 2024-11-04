package br.com.dillmann.nginxignition.application.controller.user.model

import br.com.dillmann.nginxignition.application.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.user.User.Role
import kotlinx.serialization.Serializable

@Serializable
data class UserResponse(
    val id: UuidString,
    val enabled: Boolean,
    val name: String,
    val username: String,
    val role: Role,
)
