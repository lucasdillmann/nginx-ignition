package br.com.dillmann.nginxignition.certificate.selfsigned

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.*
import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

internal class SelfSignedCertificateProvider(
    private val validator: SelfSignedCertificateValidator,
    private val factory: SelfSignedCertificateFactory,
): CertificateProvider {
    companion object {
        const val UNIQUE_ID = "SELF_SIGNED"
    }

    override val name = "Self-signed certificate"
    override val uniqueId = UNIQUE_ID
    override val dynamicFields = emptyList<CertificateProviderDynamicField>()

    override suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        validator.validate(request)
        return buildCertificate(UUID.randomUUID(), request.domainNames)
    }

    override suspend fun renew(certificate: Certificate): CertificateProvider.Output =
        buildCertificate(certificate.id, certificate.domainNames)

    private suspend fun buildCertificate(
        certificateId: UUID,
        domainNames: List<String>,
    ): CertificateProvider.Output {
        val (certificate, privateKey) = factory.build(domainNames)
        val expirationDate = certificate.notAfter.toOffsetDateTime()
        val domainModel = Certificate(
            id = certificateId,
            domainNames = domainNames,
            providerId = UNIQUE_ID,
            issuedAt = OffsetDateTime.now(),
            validFrom = certificate.notBefore.toOffsetDateTime(),
            validUntil = expirationDate,
            renewAfter = expirationDate.minusDays(1),
            privateKey = privateKey.encodeBase64(),
            publicKey = certificate.signature.encodeBase64(),
            certificationChain = emptyList(),
            answers = emptyMap(),
            metadata = null,
        )

        return CertificateProvider.Output(success = true, certificate = domainModel)
    }

    @OptIn(ExperimentalEncodingApi::class)
    private fun ByteArray.encodeBase64(): String =
        Base64.encode(this)

    private fun Date.toOffsetDateTime() =
        toInstant().atOffset(ZoneOffset.UTC)
}
