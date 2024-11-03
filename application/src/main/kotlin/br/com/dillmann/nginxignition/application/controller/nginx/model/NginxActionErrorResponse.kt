package br.com.dillmann.nginxignition.application.controller.nginx.model

import kotlinx.serialization.Serializable

@Serializable
data class NginxActionErrorResponse(
    val command: String,
    val exitCode: Int,
    val output: String,
)
