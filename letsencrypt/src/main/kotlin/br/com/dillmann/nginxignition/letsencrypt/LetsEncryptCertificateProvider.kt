package br.com.dillmann.nginxignition.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class LetsEncryptCertificateProvider(private val facade: LetsEncryptFacade): CertificateProvider {
    override val name = "Let's Encrypt"
    override val uniqueId = "LETS_ENCRYPT"
    override val dynamicFields =
        listOf(
            CertificateProviderDynamicField(
                uniqueId = "TERMS_OF_SERVICE",
                description = "I agree to the Let's Encrypt terms of service available at theirs site",
                required = true,
                type = CertificateProviderDynamicField.Type.BOOLEAN,
            ),
            CertificateProviderDynamicField(
                uniqueId = "USE_DNS_VALIDATION",
                description = "Use DNS validation",
                required = true,
                type = CertificateProviderDynamicField.Type.BOOLEAN,
            ),
            CertificateProviderDynamicField(
                uniqueId = "USER_EMAIL",
                description = "E-mail address",
                required = true,
                type = CertificateProviderDynamicField.Type.EMAIL,
            ),
        )

    override suspend fun issue(request: CertificateRequest): Certificate =
        facade.issue(request)

    override suspend fun renew(certificate: Certificate): Certificate =
        facade.renew(certificate)
}
