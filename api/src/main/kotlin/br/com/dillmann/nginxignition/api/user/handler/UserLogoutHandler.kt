package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler

internal class UserLogoutHandler(private val authorizer: Authorizer): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val principal = call.principal()
        if (principal != null)
            authorizer.revoke(principal)

        call.respond(HttpStatus.NO_CONTENT)
    }
}
