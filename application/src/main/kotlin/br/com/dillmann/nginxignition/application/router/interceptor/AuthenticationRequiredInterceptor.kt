package br.com.dillmann.nginxignition.application.router.interceptor

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.application.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.router.adapter.NettyApiCallAdapter
import org.koin.mp.KoinPlatform.getKoin

internal class AuthenticationRequiredInterceptor(private val delegate: RequestHandler): RequestHandler {
    private val authorizer = getKoin().get<RbacJwtFacade>()

    override suspend fun handle(call: ApiCall) {
        val subject = call.principal() ?: authorizer.checkCredentials(call)
        if (subject == null) {
            call.respond(HttpStatus.UNAUTHORIZED)
            return
        }

        val updatedCall = (call as NettyApiCallAdapter).copy(principal = subject)
        delegate.handle(updatedCall)
    }
}
