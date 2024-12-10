package br.com.dillmann.nginxignition.core.certificate.model

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

data class AvailableCertificateProvider(
    val id: String,
    val name: String,
    val priority: Int,
    val dynamicFields: List<DynamicField>,
)
