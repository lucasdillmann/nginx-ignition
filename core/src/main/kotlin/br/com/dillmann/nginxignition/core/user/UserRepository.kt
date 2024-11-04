package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.pagination.Page
import java.util.*

interface UserRepository {
    suspend fun save(user: User)
    suspend fun deleteById(id: UUID)
    suspend fun findById(id: UUID): User?
    suspend fun findByUsername(username: String): User?
    suspend fun findPage(pageSize: Int, pageNumber: Int): Page<User>
    suspend fun findEnabledById(id: UUID): Boolean?
    suspend fun count(): Long
}
