package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import org.koin.core.module.Module
import org.koin.dsl.binds

internal fun Module.certificateBeans() {
    single { CertificateService(get(), get(), getAll()) } binds arrayOf(
        DeleteCertificateCommand::class,
        GetCertificateCommand::class,
        IssueCertificateCommand::class,
        ListCertificateCommand::class,
        RenewCertificateCommand::class,
        GetAvailableProvidersCommand::class,
    )
    single { CertificateValidator() }
}
