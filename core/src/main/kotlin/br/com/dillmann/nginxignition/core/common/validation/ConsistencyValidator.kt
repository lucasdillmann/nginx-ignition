package br.com.dillmann.nginxignition.core.common.validation

typealias ErrorCreator = (String, String) -> Unit

internal abstract class ConsistencyValidator {
    protected suspend fun withValidationScope(scope: suspend (addError: ErrorCreator) -> Unit) {
        val errors = mutableListOf<ConsistencyException.Violation>()
        val addError: ErrorCreator = { path, message ->
            errors += ConsistencyException.Violation(path, message)
        }

        scope(addError)

        if (errors.isNotEmpty()) {
            throw ConsistencyException(errors)
        }
    }
}
