package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.user.model.UserOnboardingStatusResponse
import br.com.dillmann.nginxignition.core.user.command.GetUserCountCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus

internal class UserOnboardingStatusHandler(private val userCountCommand: GetUserCountCommand): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val output = UserOnboardingStatusResponse(finished = userCountCommand.count() > 0)
        call.respond(HttpStatus.OK, output)
    }
}
