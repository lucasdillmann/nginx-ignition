package br.com.dillmann.nginxignition.core.certificate.provider

data class CertificateRequest(
    val providerId: String,
    val hosts: List<String>,
    val answers: Map<CertificateProviderDynamicField, Any?>,
)
