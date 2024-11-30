package br.com.dillmann.nginxignition.core.certificate.command

import java.util.UUID

fun interface DeleteCertificateCommand {
    suspend fun deleteById(id: UUID)
}
