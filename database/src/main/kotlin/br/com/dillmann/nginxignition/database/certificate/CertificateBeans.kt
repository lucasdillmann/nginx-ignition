package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.certificateBeans() {
    single { CertificateDatabaseRepository() } bind CertificateRepository::class
}
