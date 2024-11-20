package br.com.dillmann.nginxignition.api.common.pagination

import kotlinx.serialization.Contextual
import kotlinx.serialization.Serializable

@Serializable
data class PageResponse<T>(
    val pageSize: Int,
    val pageNumber: Int,
    val totalItems: Long,
    val contents: List<@Contextual T>,
)
