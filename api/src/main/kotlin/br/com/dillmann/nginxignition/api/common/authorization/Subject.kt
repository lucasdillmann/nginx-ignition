package br.com.dillmann.nginxignition.api.common.authorization

import java.util.UUID

data class Subject(
    val tokenId: String = UUID.randomUUID().toString(),
    val userId: UUID,
)
