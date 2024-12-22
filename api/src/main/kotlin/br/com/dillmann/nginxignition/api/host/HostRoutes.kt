package br.com.dillmann.nginxignition.api.host

import br.com.dillmann.nginxignition.api.common.routing.*
import br.com.dillmann.nginxignition.api.host.handler.*

internal class HostRoutes(
    private val listHandler: ListHostsHandler,
    private val getByIdHandler: GetHostByIdHandler,
    private val putByIdHandler: UpdateHostByIdHandler,
    private val deleteByIdHandler: DeleteHostByIdHandler,
    private val postHandler: CreateHostHandler,
    private val toggleEnabledHandler: ToggleHostEnabledByIdHandler,
    private val getLogsHandler: GetHostLogsByIdHandler,
): RouteProvider {
    @Suppress("StringLiteralDuplication")
    override fun apiRoutes(): RouteNode =
        basePath("/api/hosts") {
            requireAuthentication {
                get(listHandler)
                post(postHandler)

                path("/{id}") {
                    get(getByIdHandler)
                    put(putByIdHandler)
                    delete(deleteByIdHandler)
                    post("/toggle-enabled", toggleEnabledHandler)
                    get("/logs/{qualifier}", getLogsHandler)
                }
            }
        }
}
