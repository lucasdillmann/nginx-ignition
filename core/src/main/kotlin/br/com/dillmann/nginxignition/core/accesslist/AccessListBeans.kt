package br.com.dillmann.nginxignition.core.accesslist

import br.com.dillmann.nginxignition.core.accesslist.command.DeleteAccessListByIdCommand
import br.com.dillmann.nginxignition.core.accesslist.command.GetAccessListByIdCommand
import br.com.dillmann.nginxignition.core.accesslist.command.ListAccessListCommand
import br.com.dillmann.nginxignition.core.accesslist.command.SaveAccessListByCommand
import org.koin.core.module.Module
import org.koin.dsl.binds

internal fun Module.accessListBeans() {
    single { AccessListService(get(), get()) } binds arrayOf(
        DeleteAccessListByIdCommand::class,
        GetAccessListByIdCommand::class,
        ListAccessListCommand::class,
        SaveAccessListByCommand::class,
    )
    single { AccessListValidator() }
}
