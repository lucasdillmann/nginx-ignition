package br.com.dillmann.nginxignition.certificate.selfsigned

import org.bouncycastle.asn1.DERSequence
import org.bouncycastle.asn1.x509.BasicConstraints
import org.bouncycastle.asn1.x509.Extension
import org.bouncycastle.asn1.x509.GeneralName
import org.bouncycastle.asn1.x509.KeyUsage
import org.bouncycastle.cert.X509CertificateHolder
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder
import org.bouncycastle.jce.provider.BouncyCastleProvider
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder
import java.math.BigInteger
import java.security.KeyPair
import java.security.KeyPairGenerator
import java.security.PrivateKey
import java.security.PublicKey
import java.security.SecureRandom
import java.security.Security
import java.time.OffsetDateTime
import java.util.*
import javax.security.auth.x500.X500Principal

class SelfSignedCertificateFactory {
    private companion object {
        private const val KEY_SIZE = 2048
    }

    private val caPrivateKey: PrivateKey
    private val caCertificate: X509CertificateHolder
    private val caPrincipal: X500Principal

    init {
        val installed = Security.getProviders().any { it is BouncyCastleProvider }
        if (!installed) {
            Security.addProvider(BouncyCastleProvider())
        }

        val caKeyPair = buildKeyPair()
        caPrivateKey = caKeyPair.private
        caPrincipal = X500Principal("CN=nginx ignition")
        caCertificate = build(caPrincipal, caKeyPair.private, caPrincipal, null, caKeyPair.public)
    }

    fun build(domainNames: List<String>): Pair<X509CertificateHolder, PrivateKey> {
        val mainDomainName = domainNames.first()
        val principal = X500Principal("CN=$mainDomainName")
        val keyPair = buildKeyPair()
        val alternativeNames =
            if (domainNames.size == 1) emptyList()
            else domainNames.drop(1).map { GeneralName(GeneralName.dNSName, it) }

        val certificate = build(
            issuer = caPrincipal,
            issuerPrivateKey = caPrivateKey,
            principal = principal,
            alternativeNames = alternativeNames,
            publicKey = keyPair.public,
        )
        return certificate to keyPair.private
    }

    private fun build(
        issuer: X500Principal,
        issuerPrivateKey: PrivateKey,
        principal: X500Principal,
        alternativeNames: List<GeneralName>?,
        publicKey: PublicKey,
    ): X509CertificateHolder {
        val baseDate = OffsetDateTime.now()
        val validFrom = Date.from(baseDate.toInstant())
        val validTo = Date.from(baseDate.plusYears(1).toInstant())

        val signer = JcaContentSignerBuilder("SHA256withRSA").build(issuerPrivateKey)
        return JcaX509v3CertificateBuilder(issuer, BigInteger.ONE, validFrom, validTo, principal, publicKey)
            .addExtension(Extension.basicConstraints, true, BasicConstraints(false))
            .addExtension(Extension.keyUsage, true, KeyUsage(KeyUsage.digitalSignature + KeyUsage.keyEncipherment))
            .also {
                if (!alternativeNames.isNullOrEmpty()) {
                    it.addExtension(
                        Extension.subjectAlternativeName,
                        true,
                        DERSequence(alternativeNames.toTypedArray()),
                    )
                }
            }
            .build(signer)
    }

    private fun buildKeyPair(): KeyPair {
        val generator = KeyPairGenerator.getInstance("RSA")
        generator.initialize(KEY_SIZE, SecureRandom())
        return generator.generateKeyPair()
    }
}
