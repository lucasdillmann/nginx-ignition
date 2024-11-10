package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest

internal class CustomCertificateProvider(
    private val validator: CustomCertificateValidator,
): CertificateProvider {
    companion object {
        const val UNIQUE_ID = "CUSTOM"
    }

    override val name = "Custom certificate"
    override val uniqueId = UNIQUE_ID
    override val dynamicFields = listOf(
        DynamicFields.CERTIFICATE_MAIN_FILE,
        DynamicFields.CERTIFICATE_CHAIN_FILE,
    )

    override suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        validator.validate(request)
        val certificateFile = request.parameters[DynamicFields.CERTIFICATE_MAIN_FILE.id] as ByteArray
        val chainFile = request.parameters[DynamicFields.CERTIFICATE_CHAIN_FILE.id] as ByteArray?
        TODO()
    }

    override suspend fun renew(certificate: Certificate) =
        CertificateProvider.Output(success = true, certificate = certificate)
}
