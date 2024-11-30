package br.com.dillmann.nginxignition.core.certificate.command

import java.util.UUID

fun interface RenewCertificateCommand {
    data class Output(
        val success: Boolean,
        val errorReason: String? = null,
    )

    suspend fun renewById(id: UUID): Output
}
