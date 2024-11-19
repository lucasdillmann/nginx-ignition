package br.com.dillmann.nginxignition.api.user.model

import kotlinx.serialization.Serializable

@Serializable
internal data class UserLoginRequest(
    val username: String,
    val password: String,
)
