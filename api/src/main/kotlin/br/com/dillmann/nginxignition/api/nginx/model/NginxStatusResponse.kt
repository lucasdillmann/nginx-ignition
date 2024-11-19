package br.com.dillmann.nginxignition.api.nginx.model

import kotlinx.serialization.Serializable

@Serializable
internal data class NginxStatusResponse(
    val running: Boolean,
)
