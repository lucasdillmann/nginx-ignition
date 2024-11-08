package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import java.util.UUID

interface IssueCertificateCommand {
    data class Output(
        val success: Boolean,
        val errorReason: String? = null,
        val certificateId: UUID? = null,
    )

    suspend fun issue(request: CertificateRequest): Output
}
