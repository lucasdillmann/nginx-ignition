package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object CustomCertificateDynamicFields {
    val PUBLIC_KEY = CertificateProviderDynamicField(
        id = "publicKey",
        description = "Certificate file (PEM encoded) with the public key",
        required = true,
        type = CertificateProviderDynamicField.Type.FILE,
    )

    val PRIVATE_KEY = CertificateProviderDynamicField(
        id = "privateKey",
        description = "Certificate file (PEM encoded) with the private key",
        required = true,
        type = CertificateProviderDynamicField.Type.FILE,
    )

    val CERTIFICATION_CHAIN = CertificateProviderDynamicField(
        id = "certificationChain",
        description = "Certification chain file (PEM encoded)",
        required = false,
        type = CertificateProviderDynamicField.Type.FILE,
    )
}
