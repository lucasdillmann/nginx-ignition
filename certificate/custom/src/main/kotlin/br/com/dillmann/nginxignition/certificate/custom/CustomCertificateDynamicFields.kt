package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object CustomCertificateDynamicFields {
    val PUBLIC_KEY = CertificateProviderDynamicField(
        id = "publicKey",
        priority = 0,
        description = "Certificate file (PEM encoded) with the public key",
        required = true,
        type = CertificateProviderDynamicField.Type.FILE,
    )

    val PRIVATE_KEY = CertificateProviderDynamicField(
        id = "privateKey",
        priority = 1,
        description = "Certificate file (PEM encoded) with the private key",
        required = true,
        type = CertificateProviderDynamicField.Type.FILE,
    )

    val CERTIFICATION_CHAIN = CertificateProviderDynamicField(
        id = "certificationChain",
        priority = 2,
        description = "Certification chain file (PEM encoded)",
        required = false,
        type = CertificateProviderDynamicField.Type.FILE,
    )
}
