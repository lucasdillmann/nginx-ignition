package br.com.dillmann.nginxignition.application.netty

import br.com.dillmann.nginxignition.application.netty.lifecycle.NettyShutdown
import br.com.dillmann.nginxignition.application.netty.lifecycle.NettyStartup
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import org.koin.core.module.Module
import org.koin.dsl.bind

fun Module.nettyBeans() {
    single { NettyContainer(get(), get()) }
    single { NettyChannelInitializer(get()) }
    single { NettyStartup(get()) } bind StartupCommand::class
    single { NettyShutdown(get()) } bind ShutdownCommand::class
}
