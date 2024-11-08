package br.com.dillmann.nginxignition.core.certificate.command

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.common.pagination.Page

interface ListCertificateCommand {
    suspend fun list(pageSize: Int, pageNumber: Int): Page<Certificate>
}
