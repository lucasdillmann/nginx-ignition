package br.com.dillmann.nginxignition.api.user.model

import kotlinx.serialization.Serializable

@Serializable
internal data class UserLoginResponse(
    val token: String,
)
