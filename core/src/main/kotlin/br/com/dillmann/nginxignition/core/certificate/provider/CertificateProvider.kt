package br.com.dillmann.nginxignition.core.certificate.provider

import br.com.dillmann.nginxignition.core.certificate.Certificate

interface CertificateProvider {
    data class Output(
        val success: Boolean,
        val errorReason: String? = null,
        val certificate: Certificate? = null,
    )

    val name: String
    val uniqueId: String
    val dynamicFields: List<CertificateProviderDynamicField>

    suspend fun issue(request: CertificateRequest): Output

    suspend fun renew(certificate: Certificate): Output
}
