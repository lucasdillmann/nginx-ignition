package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.application.router.ResponseInterceptor
import io.netty.handler.codec.http.FullHttpResponse

internal class RbacResponseInterceptor(private val facade: RbacJwtFacade): ResponseInterceptor {
    override suspend fun intercept(call: ApiCall, response: FullHttpResponse): FullHttpResponse {
        val token = call.jwtToken() ?: return response
        val credentials = facade.checkCredentials(token) ?: return response
        val updatedToken = facade.refreshToken(credentials) ?: return response

        response.headers().add("Authorization", "Bearer $updatedToken")
        return response
    }
}
