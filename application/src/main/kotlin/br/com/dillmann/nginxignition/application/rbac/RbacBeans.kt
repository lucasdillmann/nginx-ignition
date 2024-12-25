package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import br.com.dillmann.nginxignition.application.router.ResponseInterceptor
import org.koin.core.module.Module
import org.koin.dsl.bind

fun Module.rbacBeans() {
    single { RbacJwtFacade(get(), get(), get()) }
    single { RbacResponseInterceptor(get()) } bind ResponseInterceptor::class
    single { RbacJwtAuthorizer(get(), get()) } bind Authorizer::class
}
