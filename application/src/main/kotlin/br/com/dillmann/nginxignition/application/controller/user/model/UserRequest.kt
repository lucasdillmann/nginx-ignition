package br.com.dillmann.nginxignition.application.controller.user.model

import br.com.dillmann.nginxignition.core.user.User.Role
import kotlinx.serialization.Serializable

@Serializable
data class UserRequest(
    val enabled: Boolean,
    val name: String,
    val username: String,
    val password: String?,
    val role: Role,
)
