package br.com.dillmann.nginxignition.api.user.model

import kotlinx.serialization.Serializable

@Serializable
internal data class UserPasswordUpdateRequest(
    val currentPassword: String,
    val newPassword: String,
)
