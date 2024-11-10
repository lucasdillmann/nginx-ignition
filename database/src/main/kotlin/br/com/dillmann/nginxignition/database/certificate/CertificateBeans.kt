package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.certificate.lifecycle.CertificateMigrations
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.certificateBeans() {
    single { CertificateConverter() }
    single { CertificateMigrations() } bind StartupCommand::class
    single { CertificateDatabaseRepository(get()) } bind CertificateRepository::class
}
