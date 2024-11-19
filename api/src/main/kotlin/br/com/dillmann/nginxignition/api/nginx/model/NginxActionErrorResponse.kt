package br.com.dillmann.nginxignition.api.nginx.model

import kotlinx.serialization.Serializable

@Serializable
internal data class NginxActionErrorResponse(
    val command: String,
    val exitCode: Int,
    val output: String,
)
