package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import com.sun.net.httpserver.HttpExchange

internal fun interface ResponseInterceptor {
    suspend fun intercept(call: ApiCall, exchange: HttpExchange)
}
