package br.com.dillmann.nginxignition.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class LetsEncryptFacade {
    suspend fun issue(request: CertificateRequest): Certificate {
        TODO("Not yet implemented")
    }

    suspend fun renew(certificate: Certificate): Certificate {
        TODO("Not yet implemented")
    }
}
