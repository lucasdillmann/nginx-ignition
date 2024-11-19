package br.com.dillmann.nginxignition.application.api

import br.com.dillmann.nginxignition.application.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import br.com.dillmann.nginxignition.api.common.authorization.Subject

class JwtAuthorizer(
    private val jwtFacade: RbacJwtFacade,
    private val getUser: GetUserCommand,
): Authorizer {
    override suspend fun revoke(subject: Subject) {
        jwtFacade.revokeCredentials(subject.tokenId)
    }

    override suspend fun buildToken(subject: Subject): String {
        val user = getUser.getById(subject.userId)!!
        return jwtFacade.buildToken(user)
    }
}
