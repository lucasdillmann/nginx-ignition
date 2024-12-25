package br.com.dillmann.nginxignition.application.netty

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.log.logger
import io.netty.bootstrap.ServerBootstrap
import io.netty.channel.Channel
import io.netty.channel.ChannelOption
import io.netty.channel.nio.NioEventLoopGroup
import io.netty.channel.socket.nio.NioServerSocketChannel

internal class NettyContainer(
    private val configuration: ConfigurationProvider,
    private val channelInitializer: NettyChannelInitializer,
) {
    private companion object {
        private const val SO_BACKLOG_VALUE = 1024
        private val LOGGER = logger<NettyContainer>()
    }

    private lateinit var mainLoopGroup: NioEventLoopGroup
    private lateinit var workLoopGroup: NioEventLoopGroup
    private lateinit var channel: Channel

    @Synchronized
    fun start() {
        val port = configuration.get("nginx-ignition.server.port").toInt()

        mainLoopGroup = NioEventLoopGroup(1)
        workLoopGroup = NioEventLoopGroup()
        channel = ServerBootstrap()
            .option(ChannelOption.SO_BACKLOG, SO_BACKLOG_VALUE)
            .group(mainLoopGroup, workLoopGroup)
            .channel(NioServerSocketChannel::class.java)
            .childHandler(channelInitializer)
            .bind(port)
            .sync()
            .channel()

        LOGGER.info("Ready for requests on port $port")
    }

    @Synchronized
    fun stop() {
        if (::channel.isInitialized)
            channel.close()

        if (::mainLoopGroup.isInitialized)
            mainLoopGroup.shutdownGracefully()

        if (::workLoopGroup.isInitialized)
            workLoopGroup.shutdownGracefully()
    }
}
