package br.com.dillmann.nginxignition.core.common.validation

data class ConsistencyException(
    val violations: List<Violation>,
): RuntimeException() {
    data class Violation(
        val path: String,
        val value: Any?,
        val message: String,
    )
}
