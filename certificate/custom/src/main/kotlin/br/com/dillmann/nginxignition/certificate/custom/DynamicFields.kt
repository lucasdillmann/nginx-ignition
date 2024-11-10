package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object DynamicFields {
    val CERTIFICATE_MAIN_FILE = CertificateProviderDynamicField(
        id = "certificateMainFile",
        description = "Certificate file (PEM encoded) with the private key",
        required = true,
        type = CertificateProviderDynamicField.Type.FILE,
    )

    val CERTIFICATE_CHAIN_FILE = CertificateProviderDynamicField(
        id = "certificateChainFile",
        description = "Certification chain file (PEM encoded)",
        required = false,
        type = CertificateProviderDynamicField.Type.FILE,
    )
}