package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class CertificateValidator {
    suspend fun validate(request: CertificateRequest) {
        // TODO: Implement this
        // - provider is valid/exists
    }
}
