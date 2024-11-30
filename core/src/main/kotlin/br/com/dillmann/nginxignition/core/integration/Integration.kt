package br.com.dillmann.nginxignition.core.integration

data class Integration(
    val id: String,
    val enabled: Boolean,
    val parameters: Map<String, Any?>,
)
