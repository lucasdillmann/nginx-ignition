package br.com.dillmann.nginxignition.core.common.pagination

data class Page<T>(
    val pageNumber: Int,
    val pageSize: Int,
    val totalItems: Long,
    val contents: List<T>,
) {
    companion object {
        fun <T> empty() = Page<T>(pageNumber = 0, pageSize = 0, totalItems = 0, contents = emptyList())

        fun <T> of(vararg content: T) = of(content.toList())

        fun <T> of(content: List<T>) =
            Page(
                pageNumber = 0,
                pageSize = 0,
                totalItems = content.size.toLong(),
                contents = content
            )
    }

    fun <O> map(converter: (T) -> O): Page<O> =
        Page(
            pageNumber = pageNumber,
            pageSize = pageSize,
            totalItems = totalItems,
            contents = contents.map(converter),
        )

    fun <C : Comparable<C>> sortedBy(classifier: (T) -> C): Page<T> =
        copy(contents = contents.sortedBy(classifier))
}
