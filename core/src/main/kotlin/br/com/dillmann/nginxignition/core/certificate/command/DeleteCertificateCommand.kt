package br.com.dillmann.nginxignition.core.certificate.command

import java.util.UUID

fun interface DeleteCertificateCommand {
    data class Output(
        val deleted: Boolean,
        val reason: String,
    )

    suspend fun deleteById(id: UUID): Output
}
