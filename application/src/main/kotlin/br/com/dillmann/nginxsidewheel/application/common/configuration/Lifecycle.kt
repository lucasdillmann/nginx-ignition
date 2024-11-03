package br.com.dillmann.nginxsidewheel.application.common.configuration

import br.com.dillmann.nginxsidewheel.application.common.lifecycle.LifecycleManager
import io.ktor.server.application.*
import kotlinx.coroutines.runBlocking
import org.koin.ktor.ext.inject

fun Application.configureLifecycle() {
    val lifecycleManager by inject<LifecycleManager>()

    Runtime.getRuntime().addShutdownHook(Thread {
        runBlocking {
            lifecycleManager.fireShutdownEvent()
        }
    })

    runBlocking {
        lifecycleManager.fireStartupEvent()
    }
}
