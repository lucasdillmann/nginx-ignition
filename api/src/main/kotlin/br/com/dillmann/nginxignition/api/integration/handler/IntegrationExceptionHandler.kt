package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationDisabledException
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationException
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationNotConfiguredException
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationNotFoundException

suspend fun withIntegrationExceptionHandler(call: ApiCall, action: suspend () -> Unit) {
    try {
        action()
    } catch (ex: IntegrationException) {
        when (ex) {
            is IntegrationNotFoundException ->
                call.respond(HttpStatus.NOT_FOUND)

            is IntegrationDisabledException,
            is IntegrationNotConfiguredException ->
                call.respond(
                    HttpStatus.PRECONDITION_FAILED,
                    mapOf("message" to "Integration is either disabled or not yet configured"),
                )
        }
    }
}
