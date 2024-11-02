package br.com.dillmann.nginxsidewheel.core.host

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import java.util.*

interface HostRepository {
    suspend fun findById(id: UUID): Host?
    suspend fun deleteById(id: UUID)
    suspend fun save(host: Host)
    suspend fun findAll(pageSize: Int, pageNumber: Int): Page<Host>
}
