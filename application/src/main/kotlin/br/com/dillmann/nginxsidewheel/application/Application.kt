package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.configuration.*
import io.ktor.server.application.*
import io.ktor.server.netty.*

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    configureKoin()
    // TODO: Re-enable this
//    configureSecurity()
    configureHttp()
    configureRoutes()
}
