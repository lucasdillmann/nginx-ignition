package br.com.dillmann.nginxignition.core.certificate

import br.com.dillmann.nginxignition.core.certificate.command.*
import br.com.dillmann.nginxignition.core.certificate.task.CertificateAutoRenewScheduledTask
import br.com.dillmann.nginxignition.core.common.scheduler.ScheduledTask
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.koin.dsl.binds

internal fun Module.certificateBeans() {
    single { CertificateService(get(), get(), getAll(), get(), get(), get()) } binds arrayOf(
        DeleteCertificateCommand::class,
        GetCertificateCommand::class,
        IssueCertificateCommand::class,
        ListCertificateCommand::class,
        RenewCertificateCommand::class,
        GetAvailableProvidersCommand::class,
    )
    single { CertificateValidator(getAll()) }
    single { CertificateAutoRenewScheduledTask(get(), get()) } bind ScheduledTask::class
}
