package br.com.dillmann.nginxignition.certificate.selfsigned

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object SelfSignedCertificateModule {
    fun initialize() =
        module {
            single { SelfSignedCertificateValidator() }
            single { SelfSignedCertificateFactory() }
            single { SelfSignedCertificateProvider(get(), get()) } bind CertificateProvider::class
        }
}
