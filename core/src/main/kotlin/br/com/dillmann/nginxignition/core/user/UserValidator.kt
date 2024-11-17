package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

private typealias ErrorCreator = (String, String) -> Unit

internal class UserValidator(private val repository: UserRepository) {
    suspend fun validate(
        updatedState: User,
        currentState: User?,
        suppliedPassword: String?,
    ) {
        val violations = mutableListOf<ConsistencyException.Violation>()
        val addError: ErrorCreator = { path, message -> violations += ConsistencyException.Violation(path, message) }

        if (suppliedPassword.isNullOrBlank() && currentState == null)
            addError("password", "A password is required")

        if (suppliedPassword != null && suppliedPassword.length < 8)
            addError("password", "Your password should have at least 8 characters")

        val databaseUser = repository.findByUsername(updatedState.username)
        if (databaseUser != null && databaseUser.id != updatedState.id)
            addError("username", "There's already a user with the same username")

        if (updatedState.username.length < 3 || updatedState.username.isBlank())
            addError("username", "The username should have at least 3 characters")

        if (updatedState.name.length < 3 || updatedState.name.isBlank())
            addError("name", "The name should have at least 3 characters")

        if (violations.isNotEmpty())
            throw ConsistencyException(violations)
    }
}
