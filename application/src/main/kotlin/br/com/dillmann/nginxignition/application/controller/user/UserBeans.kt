package br.com.dillmann.nginxignition.application.controller.user

import br.com.dillmann.nginxignition.application.controller.user.handler.*
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import org.koin.core.module.Module
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
}
