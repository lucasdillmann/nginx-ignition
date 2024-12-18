package br.com.dillmann.nginxignition.core.accesslist

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator

internal class AccessListValidator : ConsistencyValidator() {
    @Suppress("UnusedParameter")
    suspend fun validate(accessList: AccessList) {
        withValidationScope {
            // TODO: Implement this
        }
    }
}
