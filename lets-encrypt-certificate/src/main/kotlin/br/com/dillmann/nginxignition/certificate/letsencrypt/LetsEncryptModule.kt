package br.com.dillmann.nginxignition.certificate.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.certificate.letsencrypt.acme.AcmeIssuer
import br.com.dillmann.nginxignition.certificate.letsencrypt.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.letsencrypt.dns.DnsProviderAdapter
import br.com.dillmann.nginxignition.certificate.letsencrypt.dns.Route53DnsProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object LetsEncryptModule {
    fun initialize() =
        module {
            single { AcmeIssuer() }
            single { LetsEncryptFacade(get(), get(), get()) }
            single { DnsProviderAdapter(getAll()) }
            single { LetsEncryptValidator() }
            single { Route53DnsProvider() } bind DnsProvider::class
            single { LetsEncryptCertificateProvider(get(), get()) } bind CertificateProvider::class
        }
}
