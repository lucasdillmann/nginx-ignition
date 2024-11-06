package br.com.dillmann.nginxignition.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object LetsEncryptModule {
    fun initialize() =
        module {
            single { LetsEncryptFacade() }
            single { LetsEncryptCertificateProvider(get()) } bind CertificateProvider::class
        }
}
