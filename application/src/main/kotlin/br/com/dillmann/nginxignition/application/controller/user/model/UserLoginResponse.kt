package br.com.dillmann.nginxignition.application.controller.user.model

import kotlinx.serialization.Serializable

@Serializable
data class UserLoginResponse(
    val token: String,
)
