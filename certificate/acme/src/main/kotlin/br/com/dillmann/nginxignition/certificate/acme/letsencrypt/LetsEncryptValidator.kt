package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class LetsEncryptValidator {
    suspend fun validate(request: CertificateRequest) {
        // TODO: Implement this
    }

    suspend fun validate(certificate: Certificate) {
        // TODO: Implement this
    }
}
