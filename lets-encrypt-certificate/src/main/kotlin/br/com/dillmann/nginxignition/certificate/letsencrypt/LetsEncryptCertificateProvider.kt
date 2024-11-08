package br.com.dillmann.nginxignition.certificate.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class LetsEncryptCertificateProvider(
    private val facade: LetsEncryptFacade,
    private val validator: LetsEncryptValidator,
): CertificateProvider {
    companion object {
        const val UNIQUE_ID = "LETS_ENCRYPT"
    }

    override val name = "Let's Encrypt"
    override val uniqueId = UNIQUE_ID
    override val dynamicFields = listOf(
        DynamicFields.EMAIL_ADDRESS,
        DynamicFields.AWS_ACCESS_KEY,
        DynamicFields.AWS_SECRET_KEY,
        DynamicFields.TERMS_OF_SERVICE,
    )

    override suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        validator.validate(request)
        return facade.issue(request)
    }

    override suspend fun renew(certificate: Certificate): CertificateProvider.Output {
        validator.validate(certificate)
        return facade.renew(certificate)
    }
}
