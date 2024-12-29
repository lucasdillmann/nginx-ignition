package br.com.dillmann.nginxignition.application.http

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.log.logger
import com.sun.net.httpserver.HttpServer
import java.net.InetSocketAddress

internal class HttpContainer(
    private val configuration: ConfigurationProvider,
    private val requestHandler: HttpRequestHandler,
) {
    private companion object {
        private val LOGGER = logger<HttpContainer>()
    }

    private lateinit var server: HttpServer

    fun start() {
        val port = configuration.get("nginx-ignition.server.port").toInt()
        server = HttpServer.create(InetSocketAddress(port), 0, "/", requestHandler)
        server.start()

        LOGGER.info("HTTP container started (listening for requests on port $port)")
    }

    fun stop() {
        if (!::server.isInitialized) error("Server is not running")

        val delaySeconds = configuration.get("nginx-ignition.server.shutdown-delay-seconds").toInt()
        LOGGER.info("Stopping HTTP container (waiting up to $delaySeconds seconds for the inflight requests to finish)")
        server.stop(delaySeconds)
    }
}
