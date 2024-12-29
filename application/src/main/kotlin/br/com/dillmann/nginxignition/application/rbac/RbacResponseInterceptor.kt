package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.application.router.ResponseInterceptor
import com.sun.net.httpserver.HttpExchange

internal class RbacResponseInterceptor(private val facade: RbacJwtFacade): ResponseInterceptor {
    override suspend fun intercept(call: ApiCall, exchange: HttpExchange) {
        val token = call.jwtToken() ?: return
        val credentials = facade.checkCredentials(token) ?: return
        val updatedToken = facade.refreshToken(credentials) ?: return

        exchange.responseHeaders["Authorization"] = "Bearer $updatedToken"
    }
}
