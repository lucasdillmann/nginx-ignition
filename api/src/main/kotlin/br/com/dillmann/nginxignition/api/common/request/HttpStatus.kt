package br.com.dillmann.nginxignition.api.common.request

@Suppress("MagicNumber")
enum class HttpStatus(val code: Int) {
    OK(200),
    NO_CONTENT(204),
    BAD_REQUEST(400),
    UNAUTHORIZED(401),
    FORBIDDEN(403),
    NOT_FOUND(404),
    PRECONDITION_FAILED(412),
    FAILED_DEPENDENCY(424),
    INTERNAL_SERVER_ERROR(500),
}
