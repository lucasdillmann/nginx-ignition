package br.com.dillmann.nginxignition.core.host

import br.com.dillmann.nginxignition.core.host.command.DeleteHostCommand
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import br.com.dillmann.nginxignition.core.host.command.SaveHostCommand
import org.koin.core.module.Module
import org.koin.dsl.binds

internal fun Module.hostBeans() {
    single { HostService(get(), get()) } binds arrayOf(
        SaveHostCommand::class,
        DeleteHostCommand::class,
        GetHostCommand::class,
        ListHostCommand::class,
    )
    single { HostValidator() }
}
