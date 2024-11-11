package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.certificate.commons.CertificateConstants.PRIVATE_KEY_ALGORITHM
import br.com.dillmann.nginxignition.certificate.commons.CertificateConstants.PRIVATE_KEY_FOOTER
import br.com.dillmann.nginxignition.certificate.commons.CertificateConstants.PRIVATE_KEY_HEADER
import br.com.dillmann.nginxignition.certificate.commons.CertificateConstants.PUBLIC_KEY_FOOTER
import br.com.dillmann.nginxignition.certificate.commons.CertificateConstants.PUBLIC_KEY_HEADER
import br.com.dillmann.nginxignition.certificate.custom.extensions.decodeBase64
import br.com.dillmann.nginxignition.certificate.custom.extensions.encodeBase64
import br.com.dillmann.nginxignition.certificate.custom.extensions.toOffsetDateTime
import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import java.security.KeyFactory
import java.security.PrivateKey
import java.security.cert.CertificateFactory
import java.security.cert.X509Certificate
import java.security.spec.PKCS8EncodedKeySpec
import java.time.OffsetDateTime
import java.util.*

internal class CustomCertificateProvider(
    private val validator: CustomCertificateValidator,
): CertificateProvider {
    private companion object {
        private const val UNIQUE_ID = "CUSTOM"
    }

    override val name = "Custom certificate"
    override val uniqueId = UNIQUE_ID
    override val dynamicFields = listOf(
        CustomCertificateDynamicFields.PUBLIC_KEY,
        CustomCertificateDynamicFields.PRIVATE_KEY,
        CustomCertificateDynamicFields.CERTIFICATION_CHAIN,
    )

    override suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        validator.validate(request)

        val privateKey = request.parameters[CustomCertificateDynamicFields.PRIVATE_KEY.id] as String
        val publicKey = request.parameters[CustomCertificateDynamicFields.PUBLIC_KEY.id] as String
        val chain = request.parameters[CustomCertificateDynamicFields.CERTIFICATION_CHAIN.id] as String?
        val (parsedCertificate, parsedPrivateKey) = parseCertificate(publicKey, privateKey)
        val parsedChain = parseChain(chain)

        val certificate = Certificate(
            id = UUID.randomUUID(),
            domainNames = request.domainNames,
            providerId = UNIQUE_ID,
            issuedAt = OffsetDateTime.now(),
            validUntil = parsedCertificate.notAfter.toOffsetDateTime(),
            validFrom = parsedCertificate.notBefore.toOffsetDateTime(),
            renewAfter = null,
            privateKey = parsedPrivateKey.encoded.encodeBase64(),
            publicKey = parsedCertificate.encoded.encodeBase64(),
            certificationChain = parsedChain.map { it.encoded.encodeBase64() },
            parameters = request.parameters,
            metadata = null,
        )

        return CertificateProvider.Output(
            success = true,
            certificate = certificate,
        )
    }

    private fun parseCertificate(publicKey: String, privateKey: String): Pair<X509Certificate, PrivateKey> {
        val keyFactory = KeyFactory.getInstance(PRIVATE_KEY_ALGORITHM)
        val certFactory = CertificateFactory.getInstance("X.509")

        val parsedPrivateKey = String(privateKey.decodeBase64())
            .let { parsePemBytes(it, PRIVATE_KEY_HEADER, PRIVATE_KEY_FOOTER) }
            .let { keyFactory.generatePrivate(PKCS8EncodedKeySpec(it)) }

        val parsedCertificate = String(publicKey.decodeBase64())
            .let { parsePemBytes(it, PUBLIC_KEY_HEADER, PUBLIC_KEY_FOOTER) }
            .let { certFactory.generateCertificate(it.inputStream()) as X509Certificate }

        return parsedCertificate to parsedPrivateKey
    }

    private fun parsePemBytes(contents: String, header: String, footer: String) =
        contents
            .replace("\n", "")
            .trim()
            .removePrefix(header)
            .removeSuffix(footer)
            .decodeBase64()

    private fun parseChain(chain: String?): List<X509Certificate> {
        if (chain.isNullOrBlank()) return emptyList()

        val certFactory = CertificateFactory.getInstance("X.509")
        return String(chain.decodeBase64())
            .split(PUBLIC_KEY_FOOTER)
            .filter { it.isNotBlank() }
            .map { parsePemBytes(it, PUBLIC_KEY_HEADER, PUBLIC_KEY_FOOTER) }
            .map { certFactory.generateCertificate(it.inputStream()) as X509Certificate }
    }

    override suspend fun renew(certificate: Certificate) =
        CertificateProvider.Output(success = true, certificate = certificate)
}
