package br.com.dillmann.nginxignition.core.user.lifecycle

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.user.UserService
import br.com.dillmann.nginxignition.core.user.model.SaveUserRequest
import java.util.*

internal class UserOnboarding(private val service: UserService): StartupCommand {
    private val logger = logger<UserOnboarding>()
    override val priority = 999

    override suspend fun execute() {
        val users = service.count()
        if (users > 0) return

        val firstUser = SaveUserRequest(
            id = UUID.randomUUID(),
            enabled = true,
            name = "Admin",
            username = "admin",
            password = UUID.randomUUID().toString().replace("-", ""),
            role = br.com.dillmann.nginxignition.core.user.User.Role.ADMIN,
        )
        service.save(firstUser)

        logger.info(
            "This is the first time that the application is running. You can use [admin] as username " +
            "and [${firstUser.password}] as the password (both without the brackets) to log-in.",
        )
    }
}
