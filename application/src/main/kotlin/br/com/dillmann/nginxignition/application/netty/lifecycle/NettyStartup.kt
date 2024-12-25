package br.com.dillmann.nginxignition.application.netty.lifecycle

import br.com.dillmann.nginxignition.application.netty.NettyContainer
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand

internal class NettyStartup(private val container: NettyContainer): StartupCommand {
    override suspend fun execute() {
        container.start()
    }
}
