package br.com.dillmann.nginxignition.application.netty

import br.com.dillmann.nginxignition.application.router.RequestRouter
import io.netty.channel.ChannelHandlerContext
import io.netty.channel.SimpleChannelInboundHandler
import io.netty.handler.codec.http.*
import kotlinx.coroutines.*

internal class NettyRequestHandler(
    private val router: RequestRouter,
): SimpleChannelInboundHandler<HttpObject>(false) {
    override fun channelRead0(context: ChannelHandlerContext, message: HttpObject) {
        if (message !is FullHttpRequest) return
        safeExecute(context) {
            router.route(context, message)
        }
    }

    override fun channelReadComplete(context: ChannelHandlerContext) {
        context.flush()
    }

    override fun exceptionCaught(context: ChannelHandlerContext, cause: Throwable) {
        router.route(context, cause)
    }

    private fun safeExecute(
        context: ChannelHandlerContext,
        action: suspend () -> Unit
    ) {
        // TODO: Try to integrate coroutines using native netty APIs instead of a runBlocking
        runBlocking {
            try {
                action()
            } catch (ex: Exception) {
                router.route(context, ex)
            }
        }

    }
}
