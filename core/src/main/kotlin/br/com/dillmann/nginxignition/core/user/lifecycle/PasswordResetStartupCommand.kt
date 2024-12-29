package br.com.dillmann.nginxignition.core.user.lifecycle

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.user.UserService
import java.util.UUID
import kotlin.system.exitProcess

internal class PasswordResetStartupCommand(
    private val service: UserService,
    private val configuration: ConfigurationProvider,
): StartupCommand {
    private companion object {
        private val LOGGER = logger<PasswordResetStartupCommand>()
    }

    override suspend fun execute() {
        val username =
            runCatching { configuration.get("nginx-ignition.password-reset.username") }.getOrNull() ?: return

        LOGGER.info("Starting the password reset for the username {}", username)
        val newPassword = UUID.randomUUID().toString().split("-").first()
        service.resetPassword(username, newPassword)

        LOGGER.info("Password reset completed successfully for the user {}. New password: {}", username, newPassword)
        LOGGER.warn(
            "Application was started with the password reset environment variable set and will be now terminated. " +
                "Please remove the environment variable in order to resume the normal operation/boot procedures."
        )
        exitProcess(0)
    }
}
