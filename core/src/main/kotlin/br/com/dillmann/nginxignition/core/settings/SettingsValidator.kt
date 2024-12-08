package br.com.dillmann.nginxignition.core.settings

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator
import br.com.dillmann.nginxignition.core.common.validation.ErrorCreator
import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.host.HostValidator

internal class SettingsValidator(private val bindingsValidator: HostValidator): ConsistencyValidator() {
    private companion object {
        private const val MAXIMUM_LOG_FORMAT_LENGTH = 512
        private val TIMEOUT_RANGE = 1..Int.MAX_VALUE
        private val INTERVAL_RANGE = 1..Int.MAX_VALUE
        private val LOG_LINES_RANGE = 0..10_000
        private val WORKER_PROCESSES_RANGE = 1..100
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
        }

        with(settings.logs) {
            checkLogFormatLength(accessLogsFormat, "accessLogsFormat", addError)
            checkLogFormatLength(errorLogsFormat, "errorLogsFormat", addError)
        }

        checkRange(settings.workerProcesses, WORKER_PROCESSES_RANGE, "workerProcesses", addError)
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

    private fun checkLogFormatLength(value: String?, path: String, addError: ErrorCreator) {
        if (value == null || value.length <= MAXIMUM_LOG_FORMAT_LENGTH) return
        addError("nginx.logs.$path", "Cannot have more than $MAXIMUM_LOG_FORMAT_LENGTH characters")
    }
}
