package br.com.dillmann.nginxignition.core.certificate.command

import java.util.UUID

interface RenewCertificateCommand {
    suspend fun renewById(id: UUID)
}
