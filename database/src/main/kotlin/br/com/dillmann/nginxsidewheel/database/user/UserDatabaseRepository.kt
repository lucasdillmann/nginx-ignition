package br.com.dillmann.nginxsidewheel.database.user

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.user.User
import br.com.dillmann.nginxsidewheel.core.user.UserRepository
import br.com.dillmann.nginxsidewheel.database.common.coTransaction
import br.com.dillmann.nginxsidewheel.database.user.mapping.UserTable
import org.jetbrains.exposed.sql.Op
import org.jetbrains.exposed.sql.SqlExpressionBuilder
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.upsert
import java.util.*

internal class UserDatabaseRepository(private val converter: UserConverter): UserRepository {
    override suspend fun findById(id: UUID): User? =
        findOneWhere { UserTable.id eq id }

    override suspend fun findByUsername(username: String): User? =
        findOneWhere { UserTable.username eq username }

    override suspend fun deleteById(id: UUID) {
        coTransaction {
            UserTable.deleteWhere { UserTable.id eq id }
        }
    }

    override suspend fun save(user: User) {
        coTransaction {
            UserTable.upsert { converter.apply(user, it) }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int): Page<User> =
        coTransaction {
            val totalCount = UserTable.select(UserTable.id).count()
            val users = UserTable
                .select(UserTable.columns)
                .limit(pageSize, pageSize.toLong() * pageNumber)
                .map { converter.toUser(it) }

            Page(
                contents = users,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }

    override suspend fun count(): Long =
        coTransaction {
            UserTable.select(UserTable.id).count()
        }

    private suspend fun findOneWhere(predicate: SqlExpressionBuilder.() -> Op<Boolean>): User? =
        coTransaction {
            val user = UserTable
                .select(UserTable.columns)
                .where(predicate)
                .limit(1)
                .firstOrNull()
                ?: return@coTransaction null

            converter.toUser(user)
        }
}
