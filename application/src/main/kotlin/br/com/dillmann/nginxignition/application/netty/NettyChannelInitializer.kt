package br.com.dillmann.nginxignition.application.netty

import br.com.dillmann.nginxignition.application.router.RequestRouter
import io.netty.channel.ChannelInitializer
import io.netty.channel.socket.SocketChannel
import io.netty.handler.codec.http.*

internal class NettyChannelInitializer(
    private val router: RequestRouter,
): ChannelInitializer<SocketChannel>() {
    private companion object {
        private const val MAXIMUM_CONTENT_LENGTH_MB = 1024
        private const val MAXIMUM_CONTENT_LENGTH_BYTES = MAXIMUM_CONTENT_LENGTH_MB * 1024 * 1024
    }

    override fun initChannel(channel: SocketChannel) {
        channel
            .pipeline()
            .addLast(HttpServerCodec())
            .addLast(HttpContentCompressor())
            .addLast(HttpServerExpectContinueHandler())
            .addLast(HttpObjectAggregator(MAXIMUM_CONTENT_LENGTH_BYTES))
            .addLast(NettyRequestHandler(router))
    }
}
