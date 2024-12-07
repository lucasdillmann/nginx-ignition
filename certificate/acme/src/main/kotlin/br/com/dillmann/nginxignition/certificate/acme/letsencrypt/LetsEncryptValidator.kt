package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

import br.com.dillmann.nginxignition.certificate.commons.validation.BaseCertificateValidator
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException

internal class LetsEncryptValidator: BaseCertificateValidator(
    listOf(
        LetsEncryptDynamicFields.EMAIL_ADDRESS,
        LetsEncryptDynamicFields.AWS_ACCESS_KEY,
        LetsEncryptDynamicFields.AWS_SECRET_KEY,
        LetsEncryptDynamicFields.CLOUDFLARE_API_TOKEN,
        LetsEncryptDynamicFields.DNS_PROVIDER,
    ),
) {
    override fun getDomainViolations(request: CertificateRequest): List<ConsistencyException.Violation> {
        val tosField = LetsEncryptDynamicFields.TERMS_OF_SERVICE.id
        val termsOfService = request.parameters[tosField] as? Boolean? ?: false
        if (!termsOfService) {
            return listOf(
                ConsistencyException.Violation(
                    path = "parameters.$tosField",
                    message = "You must accept the Let's Encrypt's terms of service to be able to use its certificates",
                )
            )
        }

        return emptyList()
    }
}
