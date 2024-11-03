package br.com.dillmann.nginxignition.core.common.pagination

data class Page<T>(
    val pageNumber: Int,
    val pageSize: Int,
    val totalItems: Long,
    val contents: List<T>,
)
