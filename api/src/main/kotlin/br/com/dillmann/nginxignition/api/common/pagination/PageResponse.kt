package br.com.dillmann.nginxignition.application.common.pagination

import kotlinx.serialization.Serializable

@Serializable
data class PageResponse<T>(
    val pageSize: Int,
    val pageNumber: Int,
    val totalItems: Long,
    val contents: List<T>,
)
