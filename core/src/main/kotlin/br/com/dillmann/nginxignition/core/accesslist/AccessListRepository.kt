package br.com.dillmann.nginxignition.core.accesslist

import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.UUID

interface AccessListRepository {
    suspend fun findById(id: UUID): AccessList?
    suspend fun findByName(name: String): AccessList?
    suspend fun findPage(pageNumber: Int, pageSize: Int, searchTerms: String?): Page<AccessList>
    suspend fun findAll(): List<AccessList>
    suspend fun save(accessList: AccessList)
    suspend fun deleteById(id: UUID)
}
