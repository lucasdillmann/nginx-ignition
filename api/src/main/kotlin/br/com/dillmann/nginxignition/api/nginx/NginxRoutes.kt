package br.com.dillmann.nginxignition.api.nginx

import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.basePath
import br.com.dillmann.nginxignition.api.nginx.handler.*
import br.com.dillmann.nginxignition.api.nginx.handler.NginxReloadHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStartHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStatusHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStopHandler

internal class NginxRoutes(
    private val startHandler: NginxStartHandler,
    private val stopHandler: NginxStopHandler,
    private val reloadHandler: NginxReloadHandler,
    private val statusHandler: NginxStatusHandler,
    private val getLogsHandler: GetNginxLogsHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =

    basePath("/api/nginx") {
        requireAuthentication {
            post("/start", startHandler)
            post("/stop", stopHandler)
            post("/reload", reloadHandler)
            get("/status", statusHandler)
            get("/logs", getLogsHandler)
        }
    }
}
