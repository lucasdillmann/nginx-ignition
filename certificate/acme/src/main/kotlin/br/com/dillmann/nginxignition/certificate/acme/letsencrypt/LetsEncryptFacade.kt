package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.core.certificate.provider.CertificateRequest
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider
import br.com.dillmann.nginxignition.certificate.acme.AcmeIssuer
import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProviderAdapter
import br.com.dillmann.nginxignition.certificate.acme.utils.decodeBase64
import br.com.dillmann.nginxignition.certificate.acme.utils.encodeBase64
import br.com.dillmann.nginxignition.certificate.acme.utils.toOffsetDateTime
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.shredzone.acme4j.util.KeyPairUtils
import java.security.KeyFactory
import java.security.KeyPair
import java.security.spec.PKCS8EncodedKeySpec
import java.security.spec.X509EncodedKeySpec
import java.time.Duration
import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.*

internal class LetsEncryptFacade(
    configurationProvider: ConfigurationProvider,
    private val dnsProvider: DnsProviderAdapter,
    private val acmeIssuer: AcmeIssuer,
) {
    private companion object {
        private val TIMEOUT = Duration.ofMinutes(2)
        private const val DNS_PROVIDER_ID = "ROUTE_53"
        private const val SECURITY_KEYS_ALGORITHM = "secp384r1"
    }

    private val logger = logger<LetsEncryptFacade>()
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.certificate.lets-encrypt")

    suspend fun issue(request: CertificateRequest): CertificateProvider.Output {
        val userKeys = KeyPairUtils.createKeyPair()
        val userMail = request.parameters[LetsEncryptDynamicFields.EMAIL_ADDRESS.id] as String
        val domainKeys = KeyPairUtils.createKeyPair()
        val domainNames = request.domainNames
        val productionEnvironment = configurationProvider.get("production").toBoolean()

        return issue(
            UUID.randomUUID(),
            userKeys,
            userMail,
            domainKeys,
            domainNames,
            request.parameters,
            productionEnvironment,
        )
    }

    suspend fun renew(certificate: Certificate): CertificateProvider.Output {
        val metadata = Json.decodeFromString<LetsEncryptMetadata>(certificate.metadata!!)
        val userKeys = buildKeyPair(metadata.userPrivateKey.decodeBase64(), metadata.userPublicKey.decodeBase64())
        val domainKeys = KeyPairUtils.createKeyPair()
        val parameters = certificate.parameters
        val domainNames = certificate.domainNames
        val userMail = metadata.userMail
        val productionEnvironment = metadata.productionEnvironment

        return issue(
            certificate.id,
            userKeys,
            userMail,
            domainKeys,
            domainNames,
            parameters,
            productionEnvironment,
        )
    }

    private suspend fun issue(
        certificateId: UUID,
        userKeys: KeyPair,
        userMail: String,
        domainKeys: KeyPair,
        domainNames: List<String>,
        parameters: Map<String, Any?>,
        productionEnvironment: Boolean,
    ): CertificateProvider.Output {
        val acmeContext = AcmeIssuer.Context(
            userKeys = userKeys,
            userMail = userMail,
            domainKeys = domainKeys,
            domainNames = domainNames,
            issuerUrl = certificateAuthorityUrl(productionEnvironment),
            timeout = TIMEOUT,
            createDnsRecordAction = { dnsProvider.writeChallengeRecord(DNS_PROVIDER_ID, it, parameters) },
        )
        val acmeResult = acmeIssuer.issue(acmeContext)
        if (!acmeResult.success) {
            logger.error("Certificate failed to be issued", acmeResult.exception)
            return CertificateProvider.Output(success = false, errorReason = acmeResult.exception?.message)
        }

        val issuedCertificate = acmeResult.certificate!!
        val renewInfo = issuedCertificate.renewalInfo
        val certificateChain = issuedCertificate.certificateChain
        val certificateContents = issuedCertificate.certificate
        val metadata = LetsEncryptMetadata(
            userMail = userMail,
            userPrivateKey = userKeys.private.encoded.encodeBase64(),
            userPublicKey = userKeys.public.encoded.encodeBase64(),
            productionEnvironment = productionEnvironment,
        )

        val certificate = Certificate(
            id = certificateId,
            domainNames = domainNames,
            providerId = LetsEncryptCertificateProvider.UNIQUE_ID,
            issuedAt = OffsetDateTime.now(),
            validFrom = certificateContents.notBefore.toOffsetDateTime(),
            validUntil = certificateContents.notAfter.toOffsetDateTime(),
            renewAfter = renewInfo.suggestedWindowStart.atOffset(ZoneOffset.UTC),
            privateKey = domainKeys.private.encoded.encodeBase64(),
            publicKey = certificateContents.encoded.encodeBase64(),
            certificationChain = certificateChain.map { it.encoded.encodeBase64() },
            parameters = parameters,
            metadata = Json.encodeToString(metadata),
        )
        return CertificateProvider.Output(success = true, certificate = certificate)
    }

    private fun certificateAuthorityUrl(productionEnvironment: Boolean): String {
        val environment = if (productionEnvironment) "production" else "staging"
        return "acme://letsencrypt.org/$environment"
    }

    private fun buildKeyPair(private: ByteArray, public: ByteArray): KeyPair {
        val keyFactory = KeyFactory.getInstance(SECURITY_KEYS_ALGORITHM)
        val privateKey = keyFactory.generatePrivate(PKCS8EncodedKeySpec(private))
        val publicKey = keyFactory.generatePublic(X509EncodedKeySpec(public))
        return KeyPair(publicKey, privateKey)
    }
}
