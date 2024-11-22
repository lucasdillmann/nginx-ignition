package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.host.command.*
import org.koin.core.module.Module
import org.koin.dsl.binds

internal fun Module.hostBeans() {
    single { HostService(get(), get()) } binds arrayOf(
        SaveHostCommand::class,
        DeleteHostCommand::class,
        GetHostCommand::class,
        ListHostCommand::class,
        HostExistsByIdCommand::class,
    )
    single { HostValidator(get(), get()) }
}
