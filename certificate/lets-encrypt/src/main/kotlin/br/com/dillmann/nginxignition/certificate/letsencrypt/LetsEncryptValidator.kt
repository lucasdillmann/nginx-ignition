package br.com.dillmann.nginxignition.certificate.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class LetsEncryptValidator {
    suspend fun validate(request: CertificateRequest) {
        // TODO
    }

    suspend fun validate(certificate: Certificate) {
        // TODO
    }
}
