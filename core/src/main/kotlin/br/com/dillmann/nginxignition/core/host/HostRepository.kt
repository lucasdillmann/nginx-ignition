package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

interface HostRepository {
    suspend fun findById(id: UUID): Host?
    suspend fun deleteById(id: UUID)
    suspend fun save(host: Host)
    suspend fun findPage(pageSize: Int, pageNumber: Int): Page<Host>
    suspend fun findAll(): List<Host>
}
