package br.com.dillmann.nginxignition.api.user

import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.user.handler.*
import br.com.dillmann.nginxignition.api.user.model.UserConverter
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.userBeans() {
    single { Mappers.getMapper(UserConverter::class.java) }
    single { CreateUserHandler(get(), get()) }
    single { DeleteUserByIdHandler(get()) }
    single { GetUserByIdHandler(get(), get()) }
    single { ListUsersHandler(get(), get()) }
    single { UpdateUserByIdHandler(get(), get()) }
    single { UserLoginHandler(get(), get()) }
    single { UserLogoutHandler(get()) }
    single { CurrentUserHandler(get(), get()) }
    single { UserOnboardingStatusHandler(get()) }
    single { UserOnboardingFinishHandler(get(), get(), get(), get(), get()) }
    single {
        UserRoutes(get(), get(), get(), get(), get(), get(), get(), get(), get(), get())
    } bind RouteProvider::class
}
