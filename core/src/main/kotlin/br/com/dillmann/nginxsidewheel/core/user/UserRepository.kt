package br.com.dillmann.nginxsidewheel.core.user

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import java.util.UUID

interface UserRepository {
    suspend fun save(user: User)
    suspend fun deleteById(id: UUID)
    suspend fun findById(id: UUID): User?
    suspend fun findByUsername(username: String): User?
    suspend fun findPage(pageSize: Int, pageNumber: Int): Page<User>
    suspend fun count(): Long
}
