package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.user.command.*
import br.com.dillmann.nginxignition.core.user.model.SaveUserRequest
import br.com.dillmann.nginxignition.core.user.security.UserSecurity
import java.util.*

internal class UserService(
    private val repository: UserRepository,
    private val validator: UserValidator,
    private val security: UserSecurity,
): AuthenticateUserCommand, DeleteUserCommand, GetUserCommand, ListUserCommand, SaveUserCommand,
   GetUserStatusCommand, GetUserCountCommand, UpdateUserPasswordCommand {
    override suspend fun authenticate(username: String, password: String): User? {
        val user = repository.findByUsername(username)?.takeIf { it.enabled } ?: return null
        val passwordsMatch = security.check(password, user.passwordHash, user.passwordSalt)
        return if (passwordsMatch) user else null
    }

    override suspend fun deleteById(id: UUID) {
        repository.deleteById(id)
    }

    override suspend fun getById(id: UUID): User? =
        repository.findById(id)

    override suspend fun list(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<User> =
        repository.findPage(pageSize, pageNumber, searchTerms)

    override suspend fun save(request: SaveUserRequest, currentUserId: UUID?) {
        val databaseState = repository.findById(request.id)
        val (passwordHash, passwordSalt) =
            if (request.password != null) with(security.hash(request.password)) { hash to salt }
            else if (databaseState != null) databaseState.passwordHash to databaseState.passwordSalt
            else "" to ""

        val newState = User(
            id = request.id,
            name = request.name,
            enabled = request.enabled,
            role = request.role,
            username = request.username,
            passwordHash = passwordHash,
            passwordSalt = passwordSalt,
        )
        validator.validate(newState, databaseState, request.password, currentUserId)
        repository.save(newState)
    }

    override suspend fun isEnabled(id: UUID): Boolean =
        repository.findEnabledById(id) == true

    override suspend fun count() =
        repository.count()

    override suspend fun updatePassword(userId: UUID, currentPassword: String, newPassword: String) {
        val userDetails = repository.findById(userId)!!
        validator.validatePasswordUpdate(userDetails, currentPassword, newPassword)

        val (passwordHash, passwordSalt) = security.hash(newPassword)
        val updatedUser = userDetails.copy(passwordHash = passwordHash, passwordSalt = passwordSalt)
        repository.save(updatedUser)
    }
}
