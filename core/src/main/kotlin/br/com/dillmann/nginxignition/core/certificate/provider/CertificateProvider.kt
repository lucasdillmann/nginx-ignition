package br.com.dillmann.nginxignition.core.certificate.provider

import br.com.dillmann.nginxignition.core.certificate.Certificate

interface CertificateProvider {
    val name: String
    val uniqueId: String
    val dynamicFields: List<CertificateProviderDynamicField>

    suspend fun issue(request: CertificateRequest): Certificate

    suspend fun renew(certificate: Certificate): Certificate
}
