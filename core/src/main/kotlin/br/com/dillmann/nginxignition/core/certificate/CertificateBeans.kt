package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import br.com.dillmann.nginxignition.core.certificate.lifecycle.CertificateAutoRenewStartup
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.koin.dsl.binds

internal fun Module.certificateBeans() {
    single { CertificateService(get(), get(), getAll(), get()) } binds arrayOf(
        DeleteCertificateCommand::class,
        GetCertificateCommand::class,
        IssueCertificateCommand::class,
        ListCertificateCommand::class,
        RenewCertificateCommand::class,
        GetAvailableProvidersCommand::class,
    )
    single { CertificateValidator(getAll()) }
    single { CertificateAutoRenewStartup(get(), get()) } bind StartupCommand::class
}
