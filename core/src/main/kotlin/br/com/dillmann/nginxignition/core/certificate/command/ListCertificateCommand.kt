package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.common.pagination.Page

fun interface ListCertificateCommand {
    suspend fun list(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<Certificate>
}
