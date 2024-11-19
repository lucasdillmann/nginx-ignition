package br.com.dillmann.nginxignition.api.user

import br.com.dillmann.nginxignition.api.user.handler.*
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.routes

internal class UserRoutes(
    private val listHandler: ListUsersHandler,
    private val getByIdHandler: GetUserByIdHandler,
    private val putByIdHandler: UpdateUserByIdHandler,
    private val deleteByIdHandler: DeleteUserByIdHandler,
    private val postHandler: CreateUserHandler,
    private val loginHandler: UserLoginHandler,
    private val logoutHandler: UserLogoutHandler,
    private val currentUserHandler: CurrentUserHandler,
    private val onboardingStatusHandler: UserOnboardingStatusHandler,
    private val onboardingFinishHandler: UserOnboardingFinishHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        routes("/api/users") {
            post("/login", loginHandler)
            get("/onboarding/status", onboardingStatusHandler)
            post("/onboarding/finish", onboardingFinishHandler)

            requireAuthentication {
                post("/logout", logoutHandler)
                get("/current", currentUserHandler)
            }

            requireRole(User.Role.ADMIN) {
                get(listHandler)
                get("/{id}", getByIdHandler)
                put("/{id}", putByIdHandler)
                delete("/{id}", deleteByIdHandler)
                post(postHandler)
            }
        }
}
