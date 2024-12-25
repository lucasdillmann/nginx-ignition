package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.api.common.request.ApiCall

internal suspend fun ApiCall.jwtToken() =
    headers()["authorization"]?.firstOrNull()?.substringAfter("Bearer ")
