package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import java.util.UUID

interface IssueCertificateCommand {
    suspend fun issue(request: CertificateRequest): UUID
}
