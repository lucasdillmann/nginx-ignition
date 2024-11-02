package br.com.dillmann.nginxsidewheel.application.controller.nginx.model

import kotlinx.serialization.Serializable

@Serializable
data class NginxActionErrorResponse(
    val command: String,
    val exitCode: Int,
    val output: String,
)
