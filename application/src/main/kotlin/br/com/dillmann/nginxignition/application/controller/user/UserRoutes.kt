package br.com.dillmann.nginxignition.application.controller.user

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.common.rbac.requireRole
import br.com.dillmann.nginxignition.application.common.routing.*
import br.com.dillmann.nginxignition.application.controller.user.handler.*
import br.com.dillmann.nginxignition.core.user.User
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.userRoutes() {
    val listHandler by inject<ListUsersHandler>()
    val getByIdHandler by inject<GetUserByIdHandler>()
    val putByIdHandler by inject<UpdateUserByIdHandler>()
    val deleteByIdHandler by inject<DeleteUserByIdHandler>()
    val postHandler by inject<CreateUserHandler>()
    val loginHandler by inject<UserLoginHandler>()
    val logoutHandler by inject<UserLogoutHandler>()
    val currentUserHandler by inject<CurrentUserHandler>()
    val onboardingStatusHandler by inject<UserOnboardingStatusHandler>()
    val onboardingFinishHandler by inject<UserOnboardingFinishHandler>()

    routing {
        route("/api/users") {
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
}
