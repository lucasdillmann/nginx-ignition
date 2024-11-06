package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

interface CertificateRepository {
    suspend fun findById(id: UUID): Certificate?
    suspend fun deleteById(id: UUID)
    suspend fun save(certificate: Certificate)
    suspend fun findPage(pageSize: Int, pageNumber: Int): Page<Certificate>
}
