package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator

internal class CertificateValidator(providers: List<CertificateProvider>): ConsistencyValidator() {
    private val knownProviders = providers.map { it.id }

    suspend fun validate(request: CertificateRequest) {
        withValidationScope { addError ->
            if (request.providerId !in knownProviders) {
                addError("providerId", "Invalid certificate provider. Valid values: $knownProviders.")
            }
        }
    }
}
