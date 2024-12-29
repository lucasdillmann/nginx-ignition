package br.com.dillmann.nginxignition.core.user

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.user.command.*
import br.com.dillmann.nginxignition.core.user.lifecycle.PasswordResetStartupCommand
import br.com.dillmann.nginxignition.core.user.security.UserSecurity
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
        GetUserCountCommand::class,
        UpdateUserPasswordCommand::class,
    )
    single { PasswordResetStartupCommand(get(), get()) } bind StartupCommand::class
    single { UserValidator(get(), get()) }
    single { UserSecurity(get()) }
}
