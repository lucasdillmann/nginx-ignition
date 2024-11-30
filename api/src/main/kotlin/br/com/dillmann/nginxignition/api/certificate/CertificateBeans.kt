package br.com.dillmann.nginxignition.api.certificate

import br.com.dillmann.nginxignition.api.certificate.handler.*
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.certificateBeans() {
    single { Mappers.getMapper(CertificateConverter::class.java) }
    single { ListAvailableProvidersHandler(get(), get()) }
    single { IssueCertificateHandler(get(), get()) }
    single { ListCertificatesHandler(get(), get()) }
    single { GetCertificateByIdHandler(get(), get()) }
    single { DeleteCertificateByIdHandler(get()) }
    single { RenewCertificateByIdHandler(get(), get()) }
    single { CertificateRoutes(get(), get(), get(), get(), get(), get()) } bind RouteProvider::class
}
