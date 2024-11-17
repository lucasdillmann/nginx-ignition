package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserOnboardingStatusResponse
import br.com.dillmann.nginxignition.core.user.command.GetUserCountCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class UserOnboardingStatusHandler(private val userCountCommand: GetUserCountCommand): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val output = UserOnboardingStatusResponse(finished = userCountCommand.count() > 0)
        call.respond(HttpStatusCode.OK, output)
    }
}
