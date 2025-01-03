package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.Certificate
import java.util.UUID

fun interface GetCertificateCommand {
    suspend fun getById(id: UUID): Certificate?
}
