package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

class CertificateDatabaseRepository: CertificateRepository {
    override suspend fun findById(id: UUID): Certificate? {
        TODO("Not yet implemented")
    }

    override suspend fun deleteById(id: UUID) {
        TODO("Not yet implemented")
    }

    override suspend fun save(certificate: Certificate) {
        TODO("Not yet implemented")
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int): Page<Certificate> {
        TODO("Not yet implemented")
    }
}
