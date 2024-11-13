package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

internal class CertificateValidator(providers: List<CertificateProvider>) {
    private val knownProviders = providers.map { it.uniqueId }

    fun validate(request: CertificateRequest) {
        if (request.providerId !in knownProviders) {
            val violation = ConsistencyException.Violation(
                path = "providerId",
                message = "Invalid certificate provider. Valid values: $knownProviders."
            )

            throw ConsistencyException(listOf(violation))
        }
    }
}
