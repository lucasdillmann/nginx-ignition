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
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        routes("/api/hosts") {
            requireAuthentication {
                get(listHandler)
                post(postHandler)
                get("/{id}", getByIdHandler)
                put("/{id}", putByIdHandler)
                delete("/{id}", deleteByIdHandler)
                post("/{id}/toggle-enabled", toggleEnabledHandler)
                get("/{id}/access-logs") { TODO() }
                get("/{id}/error-logs") { TODO() }
            }
        }
}
