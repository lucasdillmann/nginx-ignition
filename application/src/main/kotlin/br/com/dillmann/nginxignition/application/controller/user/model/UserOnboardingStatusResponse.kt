package br.com.dillmann.nginxignition.application.controller.user.model

import kotlinx.serialization.Serializable

@Serializable
data class UserOnboardingStatusResponse(
    val finished: Boolean,
)
