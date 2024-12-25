package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import io.netty.handler.codec.http.FullHttpResponse

fun interface ResponseInterceptor {
    suspend fun intercept(call: ApiCall, response: FullHttpResponse): FullHttpResponse
}
