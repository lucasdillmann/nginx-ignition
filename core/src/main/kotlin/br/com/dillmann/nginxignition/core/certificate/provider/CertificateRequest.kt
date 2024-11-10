package br.com.dillmann.nginxignition.core.certificate.provider

data class CertificateRequest(
    val providerId: String,
    val domainNames: List<String>,
    val parameters: Map<String, Any?>,
)
