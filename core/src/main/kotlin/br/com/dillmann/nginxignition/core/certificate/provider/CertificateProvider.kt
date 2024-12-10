package br.com.dillmann.nginxignition.core.certificate.provider

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

interface CertificateProvider {
    data class Output(
        val success: Boolean,
        val errorReason: String? = null,
        val certificate: Certificate? = null,
    )

    val id: String
    val name: String
    val dynamicFields: List<DynamicField>
    val priority: Int

    suspend fun issue(request: CertificateRequest): Output

    suspend fun renew(certificate: Certificate): Output
}
