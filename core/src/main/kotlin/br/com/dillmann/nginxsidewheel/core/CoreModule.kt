package br.com.dillmann.nginxsidewheel.core

import br.com.dillmann.nginxsidewheel.core.host.HostService
import br.com.dillmann.nginxsidewheel.core.host.command.DeleteHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.GetHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.ListHostCommand
import br.com.dillmann.nginxsidewheel.core.host.command.SaveHostCommand
import org.koin.dsl.module

object CoreModule {
    fun initialize() =
        module {
            single { HostService(get()) }
            single<SaveHostCommand> { get<HostService>() }
            single<DeleteHostCommand> { get<HostService>() }
            single<GetHostCommand> { get<HostService>() }
            single<ListHostCommand> { get<HostService>() }
        }
}
