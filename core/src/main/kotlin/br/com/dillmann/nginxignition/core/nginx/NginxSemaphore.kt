package br.com.dillmann.nginxignition.core.nginx

import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock

internal class NginxSemaphore {
    enum class State {
        RUNNING,
        STOPPED,
    }

    private val lock = Mutex()
    var currentState = State.STOPPED
        private set

    suspend fun changeState(newState: State, action: suspend () -> Unit) {
        lock.withLock {
            action()
            currentState = newState
        }
    }
}
