package br.com.dillmann.nginxignition.api.common.routing

import br.com.dillmann.nginxignition.core.user.User

class RoleRequiredRouteNode(val role: User.Role): CompositeRouteNode()
