package br.com.dillmann.nginxsidewheel.core.user

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.core.user.command.*
import br.com.dillmann.nginxsidewheel.core.user.security.UserSecurity
import br.com.dillmann.nginxsidewheel.core.user.startup.UserOnboarding
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.koin.dsl.binds

internal fun Module.userBeans() {
    single { UserService(get(), get(), get()) } binds arrayOf(
        AuthenticateUserCommand::class,
        DeleteUserCommand::class,
        GetUserCommand::class,
        ListUserCommand::class,
        SaveUserCommand::class,
    )
    single { UserValidator(get()) }
    single { UserSecurity(get()) }
    single { UserOnboarding(get()) } bind StartupCommand::class
}
