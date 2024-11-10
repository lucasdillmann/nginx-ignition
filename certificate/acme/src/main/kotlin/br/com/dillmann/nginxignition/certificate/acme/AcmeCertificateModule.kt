package br.com.dillmann.nginxignition.certificate.acme

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProvider
import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.dns.DnsProviderAdapter
import br.com.dillmann.nginxignition.certificate.acme.dns.Route53DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptCertificateProvider
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptFacade
import br.com.dillmann.nginxignition.certificate.acme.letsencrypt.LetsEncryptValidator
import org.koin.dsl.bind
import org.koin.dsl.module

object AcmeCertificateModule {
    fun initialize() =
        module {
            single { AcmeIssuer() }
            single { DnsProviderAdapter(getAll()) }
            single { Route53DnsProvider() } bind DnsProvider::class
            single { LetsEncryptFacade(get(), get(), get()) }
            single { LetsEncryptValidator() }
            single { LetsEncryptCertificateProvider(get(), get()) } bind CertificateProvider::class
        }
}
