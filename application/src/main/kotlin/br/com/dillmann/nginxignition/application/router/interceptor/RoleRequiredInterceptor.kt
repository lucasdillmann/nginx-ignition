package br.com.dillmann.nginxignition.application.router.interceptor

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.application.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.router.adapter.NettyApiCallAdapter
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import org.koin.mp.KoinPlatform.getKoin

internal class RoleRequiredInterceptor(
    private val role: User.Role,
    private val delegate: RequestHandler,
): RequestHandler {
    private val authorizer = getKoin().get<RbacJwtFacade>()
    private val getUserCommand = getKoin().get<GetUserCommand>()

    override suspend fun handle(call: ApiCall) {
        val subject = call.principal() ?: authorizer.checkCredentials(call)
        val user = subject?.userId?.let { getUserCommand.getById(it) }

        if (subject == null || user == null) {
            call.respond(HttpStatus.UNAUTHORIZED)
            return
        }

        if (user.role != role) {
            call.respond(HttpStatus.FORBIDDEN)
            return
        }

        val updatedCall = (call as NettyApiCallAdapter).copy(principal = subject)
        delegate.handle(updatedCall)
    }
}
