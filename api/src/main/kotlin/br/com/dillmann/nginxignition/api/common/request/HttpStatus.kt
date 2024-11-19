package br.com.dillmann.nginxignition.api.common.request

enum class HttpStatus(val code: Int) {
    OK(200),
    NO_CONTENT(204),
    NOT_FOUND(404),
    BAD_REQUEST(400),
    FORBIDDEN(403),
    FAILED_DEPENDENCY(424),
}
