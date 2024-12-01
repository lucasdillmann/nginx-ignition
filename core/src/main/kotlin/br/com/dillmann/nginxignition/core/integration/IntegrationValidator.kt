package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicFields
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

internal class IntegrationValidator {
    fun validate(dynamicFields: List<DynamicField>, parameters: Map<String, Any?>) {
        val violations = DynamicFields.validate(dynamicFields, parameters)
        if (violations.isNotEmpty())
            throw ConsistencyException(violations)
    }
}
