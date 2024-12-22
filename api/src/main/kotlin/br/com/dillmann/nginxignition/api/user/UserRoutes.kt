package br.com.dillmann.nginxignition.api.user

import br.com.dillmann.nginxignition.api.user.handler.*
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.basePath

@Suppress("LongParameterList")
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
    private val updatePasswordHandler: UpdateCurrentUserPasswordHandler,
): RouteProvider {
    @Suppress("StringLiteralDuplication")
    override fun apiRoutes(): RouteNode =
        basePath("/api/users") {
            post("/login", loginHandler)

            path("/onboarding") {
                get("/status", onboardingStatusHandler)
                post("/finish", onboardingFinishHandler)
            }

            requireAuthentication {
                post("/logout", logoutHandler)

                path("/current") {
                    get(currentUserHandler)
                    post("/update-password", updatePasswordHandler)
                }
            }

            requireRole(User.Role.ADMIN) {
                get(listHandler)
                post(postHandler)

                path("/{id}") {
                    get(getByIdHandler)
                    put(putByIdHandler)
                    delete(deleteByIdHandler)
                }
            }
        }
}
