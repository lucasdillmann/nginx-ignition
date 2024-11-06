package br.com.dillmann.nginxignition.core.certificate.model

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

data class AvailableCertificateProvider(
    val uniqueId: String,
    val name: String,
    val dynamicFields: List<CertificateProviderDynamicField>,
)
