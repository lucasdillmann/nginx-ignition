package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import br.com.dillmann.nginxignition.core.user.security.UserSecurity
import java.util.UUID

private typealias ErrorCreator = (String, String) -> Unit

internal class UserValidator(
    private val repository: UserRepository,
    private val security: UserSecurity,
) {
    private companion object {
        private const val MINIMUM_PASSWORD_LENGTH = 8
        private const val MINIMUM_USERNAME_LENGTH = 3
    }

    suspend fun validate(
        updatedState: User,
        currentState: User?,
        suppliedPassword: String?,
        currentUserId: UUID?,
    ) {
        val violations = mutableListOf<ConsistencyException.Violation>()
        val addError: ErrorCreator = { path, message -> violations += ConsistencyException.Violation(path, message) }

        if (!updatedState.enabled && currentState != null && currentState.id == currentUserId)
            addError("enabled", "You cannot disable your own user")

        if (suppliedPassword.isNullOrBlank() && currentState == null)
            addError("password", "A password is required")

        if (suppliedPassword != null && suppliedPassword.length < MINIMUM_PASSWORD_LENGTH)
            addError("password", "Password should have at least $MINIMUM_PASSWORD_LENGTH characters")

        val databaseUser = repository.findByUsername(updatedState.username)
        if (databaseUser != null && databaseUser.id != updatedState.id)
            addError("username", "There's already a user with the same username")

        if (updatedState.username.length < MINIMUM_USERNAME_LENGTH || updatedState.username.isBlank())
            addError("username", "The username should have at least $MINIMUM_USERNAME_LENGTH characters")

        if (updatedState.name.length < MINIMUM_USERNAME_LENGTH || updatedState.name.isBlank())
            addError("name", "The name should have at least $MINIMUM_USERNAME_LENGTH characters")

        if (violations.isNotEmpty())
            throw ConsistencyException(violations)
    }

    fun validatePasswordUpdate(
        user: User,
        currentPassword: String,
        newPassword: String,
    ) {
        val violations = mutableListOf<ConsistencyException.Violation>()
        val addError: ErrorCreator = { path, message -> violations += ConsistencyException.Violation(path, message) }

        if (!security.check(currentPassword, user.passwordHash, user.passwordSalt))
            addError("currentPassword", "Not your current password")

        if (newPassword.length < MINIMUM_PASSWORD_LENGTH)
            addError("newPassword", "Your new password should have at least 8 characters")

        if (violations.isNotEmpty())
            throw ConsistencyException(violations)
    }
}
