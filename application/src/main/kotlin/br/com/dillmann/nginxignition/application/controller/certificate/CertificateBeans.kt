package br.com.dillmann.nginxignition.application.controller.certificate

import br.com.dillmann.nginxignition.application.controller.certificate.handler.IssueCertificateHandler
import br.com.dillmann.nginxignition.application.controller.certificate.handler.ListAvailableProvidersHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import org.koin.core.module.Module
import org.mapstruct.factory.Mappers

internal fun Module.certificateBeans() {
    single { Mappers.getMapper(CertificateConverter::class.java) }
    single { ListAvailableProvidersHandler(get(), get()) }
    single { IssueCertificateHandler(get(), get()) }
}
