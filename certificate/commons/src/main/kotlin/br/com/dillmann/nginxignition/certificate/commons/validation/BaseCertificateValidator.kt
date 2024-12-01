package br.com.dillmann.nginxignition.certificate.commons.validation

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.GlobalConstants.TLD_PATTERN
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicFields
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

abstract class BaseCertificateValidator(
    private val dynamicFields: List<DynamicField>,
) {
    fun validate(request: CertificateRequest) {
        val violations = validateBaseFields(request) +
            DynamicFields.validate(dynamicFields, request.parameters) +
            getDomainViolations(request)
        if (violations.isNotEmpty()) {
            throw ConsistencyException(violations)
        }
    }

    open fun getDomainViolations(request: CertificateRequest): List<ConsistencyException.Violation> =
        emptyList()

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
