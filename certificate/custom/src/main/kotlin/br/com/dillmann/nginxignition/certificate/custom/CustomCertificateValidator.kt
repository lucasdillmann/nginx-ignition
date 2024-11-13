package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.certificate.commons.validation.BaseCertificateValidator

internal class CustomCertificateValidator: BaseCertificateValidator(
    listOf(
        CustomCertificateDynamicFields.PUBLIC_KEY,
        CustomCertificateDynamicFields.PRIVATE_KEY,
        CustomCertificateDynamicFields.CERTIFICATION_CHAIN,
    ),
)
