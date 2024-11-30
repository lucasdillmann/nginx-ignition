package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object CustomCertificateDynamicFields {
    val PUBLIC_KEY = DynamicField(
        id = "publicKey",
        priority = 0,
        description = "Certificate file (PEM encoded) with the public key",
        required = true,
        sensitive = true,
        type = DynamicField.Type.FILE,
    )

    val PRIVATE_KEY = DynamicField(
        id = "privateKey",
        priority = 1,
        description = "Certificate file (PEM encoded) with the private key",
        required = true,
        sensitive = true,
        type = DynamicField.Type.FILE,
    )

    val CERTIFICATION_CHAIN = DynamicField(
        id = "certificationChain",
        priority = 2,
        description = "Certification chain file (PEM encoded)",
        required = false,
        sensitive = true,
        type = DynamicField.Type.FILE,
    )
}
