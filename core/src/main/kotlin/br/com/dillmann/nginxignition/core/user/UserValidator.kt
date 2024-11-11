package br.com.dillmann.nginxignition.core.user

internal class UserValidator(private val repository: UserRepository) {
    suspend fun validate(updatedState: User, currentState: User?) {
        // TODO: Implement this
        //  - If password is empty, currentState cannot be null (currentState = null means new user)
    }
}
