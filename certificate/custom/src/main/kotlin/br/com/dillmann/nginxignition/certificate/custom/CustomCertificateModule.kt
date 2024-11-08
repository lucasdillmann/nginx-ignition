package br.com.dillmann.nginxignition.certificate.custom

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object CustomCertificateModule {
    fun initialize() =
        module {
            single { CustomCertificateValidator() }
            single { CustomCertificateProvider(get()) } bind CertificateProvider::class
        }
}
