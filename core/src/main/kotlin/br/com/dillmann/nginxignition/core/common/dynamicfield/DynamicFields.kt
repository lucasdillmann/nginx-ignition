package br.com.dillmann.nginxignition.core.common.dynamicfield

import br.com.dillmann.nginxignition.core.common.GlobalConstants.EMAIL_PATTERN
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField.Type.*
import br.com.dillmann.nginxignition.core.common.extensions.decodeBase64
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

object DynamicFields {
    fun removeSensitiveParameters(
        dynamicFields: List<DynamicField>,
        parameters: Map<String, Any?>,
    ): Map<String, Any?> {
        val fieldsToRemove = dynamicFields.filter { it.sensitive }.map { it.id }.toSet()
        return parameters - fieldsToRemove
    }

    fun validate(
        dynamicFields: List<DynamicField>,
        parameters: Map<String, Any?>,
    ): List<ConsistencyException.Violation> {
        val violations = mutableListOf<ConsistencyException.Violation>()
        dynamicFields.forEach { field ->
            val value = parameters[field.id]
            if (value == null && field.required) {
                violations += ConsistencyException.Violation(
                    path = "parameters.${field.id}",
                    message = "A value is required",
                )
            }

            if (value != null) {
                val enumOptions = field.enumOptions.map { it.id }
                val incompatibleMessage = resolveErrorMessage(field, value, enumOptions)

                if (incompatibleMessage != null) {
                    violations += ConsistencyException.Violation(
                        path = "parameters.${field.id}",
                        message = incompatibleMessage,
                    )
                }
            }
        }

        return violations
    }

    private fun resolveErrorMessage(
        field: DynamicField,
        value: Any,
        enumOptions: List<String>
    ): String? =
        when {
            field.type in listOf(ENUM, SINGLE_LINE_TEXT, MULTI_LINE_TEXT) && value !is String ->
                "A text value is expected"

            field.type == ENUM && value !in enumOptions ->
                "Not a recognized option. Valid values: $enumOptions."

            field.type == FILE && !canDecodeFile(value) ->
                "A file is expected, encoded in a Base64 String"

            field.type == BOOLEAN && value !is Boolean ->
                "A boolean value is expected"

            field.type == EMAIL && !isAnEmail(value) ->
                "A email is expected"

            else -> null
        }

    private fun canDecodeFile(value: Any): Boolean =
        runCatching { (value as String).decodeBase64() }.getOrNull() != null

    private fun isAnEmail(value: Any): Boolean =
        value is String && EMAIL_PATTERN.matcher(value).find()
}
