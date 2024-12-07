package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

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
    override val id = UNIQUE_ID
    override val dynamicFields = listOf(
        LetsEncryptDynamicFields.EMAIL_ADDRESS,
        LetsEncryptDynamicFields.AWS_ACCESS_KEY,
        LetsEncryptDynamicFields.AWS_SECRET_KEY,
        LetsEncryptDynamicFields.CLOUDFLARE_API_TOKEN,
        LetsEncryptDynamicFields.TERMS_OF_SERVICE,
        LetsEncryptDynamicFields.DNS_PROVIDER,
    )

    override suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        validator.validate(request)
        return facade.issue(request)
    }

    override suspend fun renew(certificate: Certificate): CertificateProvider.Output {
        return facade.renew(certificate)
    }
}
