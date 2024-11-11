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
import java.security.SecureRandom
import java.security.Security
import java.time.OffsetDateTime
import java.util.*
import javax.security.auth.x500.X500Principal


class SelfSignedCertificateFactory {
    init {
        val installed = Security.getProviders().any { it is BouncyCastleProvider }
        if (!installed) {
            Security.addProvider(BouncyCastleProvider())
        }
    }

    fun build(domainNames: List<String>): Pair<X509CertificateHolder, PrivateKey> {
        val mainDomainName = domainNames.first()
        val principal = X500Principal("CN=$mainDomainName")
        val keyPair = buildKeyPair()
        val alternativeNames =
            if (domainNames.size == 1) emptyList()
            else domainNames.drop(1).map { GeneralName(GeneralName.dNSName, it) }

        val baseDate = OffsetDateTime.now()
        val validFrom = Date.from(baseDate.toInstant())
        val validTo = Date.from(baseDate.plusYears(1).toInstant())

        val signer = JcaContentSignerBuilder("SHA256withRSA").build(keyPair.private)
        val certificate =
            JcaX509v3CertificateBuilder(principal, BigInteger.ONE, validFrom, validTo, principal, keyPair.public)
                .addExtension(Extension.basicConstraints, true, BasicConstraints(false))
                .addExtension(Extension.keyUsage, true, KeyUsage(KeyUsage.digitalSignature + KeyUsage.keyEncipherment))
                .addExtension(Extension.subjectAlternativeName, true, DERSequence(alternativeNames.toTypedArray()))
                .build(signer)

        return certificate to keyPair.private
    }

    private fun buildKeyPair(): KeyPair {
        val generator = KeyPairGenerator.getInstance("RSA")
        generator.initialize(2048, SecureRandom())
        return generator.generateKeyPair()
    }
}
