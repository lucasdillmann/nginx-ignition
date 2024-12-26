package br.com.dillmann.nginxignition.application.frontend

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import org.apache.tika.Tika
import java.net.URL
import java.net.URLDecoder

internal class FrontendRequestHandler: RequestHandler {
    private companion object {
        private const val INDEX_FILE_PATH = "/index.html"
        private val INDEX_FILE = FrontendFileLoader.load(INDEX_FILE_PATH)!!
        private val BASE_PATH = INDEX_FILE.path.removeSuffix(INDEX_FILE_PATH)
    }

    override suspend fun handle(call: ApiCall) {
        val file = resolveFile(call.uri())
        if (file == null) {
            call.respond(HttpStatus.NOT_FOUND)
            return
        }

        val contents = file.readBytes()
        val contentType = Tika().detect(file)
        val headers = mapOf(
            "content-type" to contentType,
            "content-length" to contents.size.toString(),
        )

        call.respondRaw(HttpStatus.OK, headers, contents)
    }

    private fun resolveFile(uri: String): URL? {
        val decodedUri = URLDecoder.decode(uri, Charsets.UTF_8)
        val file = FrontendFileLoader.load(decodedUri)?.takeIf { !it.path.endsWith("/") } ?: INDEX_FILE

        if (!file.path.startsWith(BASE_PATH) || file.path.contains(".."))
            return null

        return file
    }
}