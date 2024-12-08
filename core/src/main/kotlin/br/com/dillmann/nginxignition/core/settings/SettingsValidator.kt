package br.com.dillmann.nginxignition.core.settings

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator
import br.com.dillmann.nginxignition.core.common.validation.ErrorCreator
import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.host.HostValidator

internal class SettingsValidator(private val bindingsValidator: HostValidator): ConsistencyValidator() {
    private companion object {
        private const val MAXIMUM_DEFAULT_CONTENT_TYPE_LENGTH = 128
        private val TIMEOUT_RANGE = 1..Int.MAX_VALUE
        private val INTERVAL_RANGE = 1..Int.MAX_VALUE
        private val LOG_LINES_RANGE = 0..10_000
        private val WORKER_PROCESSES_RANGE = 1..100
        private val WORKER_CONNECTIONS_RANGE = 32..4096
        private val MAXIMUM_BODY_SIZE_RANGE = 1..Int.MAX_VALUE
    }

    suspend fun validate(settings: Settings) {
        withValidationScope { addError ->
            validateNginx(settings.nginx, addError)
            validateLogRotation(settings.logRotation, addError)
            validateCertificateAutoRenew(settings.certificateAutoRenew, addError)
            validateGlobalBindings(settings.globalBindings, addError)
        }
    }

    private fun validateNginx(settings: Settings.NginxSettings, addError: ErrorCreator) {
        with(settings.timeouts) {
            checkRange(read, TIMEOUT_RANGE, "nginx.timeouts.read", addError)
            checkRange(send, TIMEOUT_RANGE, "nginx.timeouts.send", addError)
            checkRange(connect, TIMEOUT_RANGE, "nginx.timeouts.connect", addError)
            checkRange(keepalive, TIMEOUT_RANGE, "nginx.timeouts.keepalive", addError)
        }

        checkRange(settings.workerProcesses, WORKER_PROCESSES_RANGE, "nginx.workerProcesses", addError)
        checkRange(settings.workerConnections, WORKER_CONNECTIONS_RANGE, "nginx.workerConnections", addError)
        checkRange(settings.maximumBodySizeMb, MAXIMUM_BODY_SIZE_RANGE, "nginx.maximumBodySizeMb", addError)

        if (settings.defaultContentType.isBlank())
            addError("nginx.defaultContentType", "A value is required")

        if (settings.defaultContentType.length > MAXIMUM_DEFAULT_CONTENT_TYPE_LENGTH)
            addError("nginx.defaultContentType", "Cannot have more than 128 characters")
    }

    private fun validateLogRotation(settings: Settings.LogRotation, addError: ErrorCreator) {
        with(settings) {
            checkRange(intervalUnitCount, INTERVAL_RANGE, "logRotation.intervalUnitCount", addError)
            checkRange(maximumLines, LOG_LINES_RANGE, "logRotation.maximumLines", addError)
        }
    }

    private fun validateCertificateAutoRenew(settings: Settings.CertificateAutoRenew, addError: ErrorCreator) {
        checkRange(settings.intervalUnitCount, INTERVAL_RANGE, "certificateAutoRenew.intervalUnitCount", addError)
    }

    private suspend fun validateGlobalBindings(settings: List<Host.Binding>, addError: ErrorCreator) {
        settings.forEachIndexed { index, binding ->
            bindingsValidator.validateBinding("globalBindings", binding, index, addError)
        }
    }

    private fun checkRange(value: Int, range: IntRange, path: String, addError: ErrorCreator) {
        if (value !in range)
            addError(path, "Must be between ${range.first} and ${range.last}")
    }
}
