package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.user.command.*
import br.com.dillmann.nginxignition.core.user.security.UserSecurity
import br.com.dillmann.nginxignition.core.user.lifecycle.UserOnboarding
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.koin.dsl.binds

internal fun Module.userBeans() {
    single { UserService(get(), get(), get()) } binds arrayOf(
        AuthenticateUserCommand::class,
        DeleteUserCommand::class,
        GetUserCommand::class,
        GetUserStatusCommand::class,
        ListUserCommand::class,
        SaveUserCommand::class,
    )
    single { UserValidator(get()) }
    single { UserSecurity(get()) }
    single { UserOnboarding(get()) } bind StartupCommand::class
}
