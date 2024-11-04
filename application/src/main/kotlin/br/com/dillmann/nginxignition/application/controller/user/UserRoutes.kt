package br.com.dillmann.nginxignition.application.controller.user

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.common.rbac.requireRole
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

    routing {
        route("/api/users") {
            post("/login") { loginHandler.handle(call) }

            requireAuthentication {
                post("/logout") { logoutHandler.handle(call) }
            }

            requireRole(User.Role.ADMIN) {
                get { listHandler.handle(call) }
                get("/{id}") { getByIdHandler.handle(call) }
                put("/{id}") { putByIdHandler.handle(call) }
                delete("/{id}") { deleteByIdHandler.handle(call) }
                post { postHandler.handle(call) }
            }
        }
    }
}
