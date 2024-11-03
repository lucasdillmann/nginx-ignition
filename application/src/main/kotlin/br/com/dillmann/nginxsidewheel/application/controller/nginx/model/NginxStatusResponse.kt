package br.com.dillmann.nginxsidewheel.application.controller.nginx.model

import kotlinx.serialization.Serializable

@Serializable
data class NginxStatusResponse(
    val running: Boolean,
)
