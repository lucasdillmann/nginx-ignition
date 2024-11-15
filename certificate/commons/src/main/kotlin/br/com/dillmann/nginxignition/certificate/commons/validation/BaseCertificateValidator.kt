package br.com.dillmann.nginxignition.certificate.commons.validation

import br.com.dillmann.nginxignition.certificate.commons.extensions.decodeBase64
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField.Type.*
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.GlobalConstants.EMAIL_PATTERN
import br.com.dillmann.nginxignition.core.common.GlobalConstants.TLD_PATTERN
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

abstract class BaseCertificateValidator(
    private val dynamicFields: List<CertificateProviderDynamicField>,
) {
    fun validate(request: CertificateRequest) {
        val violations = validateBaseFields(request) + validateDynamicFields(request) + getDomainViolations(request)
        if (violations.isNotEmpty()) {
            throw ConsistencyException(violations)
        }
    }

    open fun getDomainViolations(request: CertificateRequest): List<ConsistencyException.Violation> =
        emptyList()

    private fun validateDynamicFields(request: CertificateRequest): List<ConsistencyException.Violation> {
        val violations = mutableListOf<ConsistencyException.Violation>()
        dynamicFields.forEach { field ->
            val value = request.parameters[field.id]
            if (value == null && field.required) {
                violations += ConsistencyException.Violation(
                    path = "parameters.${field.id}",
                    message = "A value is required",
                )
            }

            if (value != null) {
                val enumOptions = field.enumOptions.map { it.id }
                val incompatibleMessage =
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

    private fun canDecodeFile(value: Any): Boolean =
        runCatching { (value as String).decodeBase64() }.getOrNull() != null

    private fun isAnEmail(value: Any): Boolean =
        value is String && EMAIL_PATTERN.matcher(value).find()

    private fun validateBaseFields(request: CertificateRequest): List<ConsistencyException.Violation> {
        val violations = mutableListOf<ConsistencyException.Violation>()
        if (request.domainNames.isEmpty()) {
            violations += ConsistencyException.Violation(
                path = "domainNames",
                message = "At least one domain name must be informed",
            )
        }

        request.domainNames.forEachIndexed { index, domainName ->
            if (!TLD_PATTERN.matcher(domainName).find()) {
                violations += ConsistencyException.Violation(
                    path = "domainNames[$index]",
                    message = "Value is not a valid domain name",
                )
            }
        }

        return violations
    }
}
