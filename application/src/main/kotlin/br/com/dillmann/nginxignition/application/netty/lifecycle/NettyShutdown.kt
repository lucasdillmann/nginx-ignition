package br.com.dillmann.nginxignition.application.netty.lifecycle

import br.com.dillmann.nginxignition.application.netty.NettyContainer
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand

internal class NettyShutdown(private val container: NettyContainer): ShutdownCommand {
    override suspend fun execute() {
        container.stop()
    }
}
