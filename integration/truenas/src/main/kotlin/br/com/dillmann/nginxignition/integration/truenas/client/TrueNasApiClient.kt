package br.com.dillmann.nginxignition.integration.truenas.client

import kotlinx.serialization.json.Json
import okhttp3.Credentials
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response

internal class TrueNasApiClient(
    private val baseUrl: String,
    private val username: String,
    private val password: String,
) {
    private val jsonParser = Json { ignoreUnknownKeys = true }
    private val client = OkHttpClient
        .Builder()
        .addInterceptor { chain ->
            val credentials = Credentials.basic(username, password)
            val updatedRequest = chain.request().newBuilder().addHeader("Authorization", credentials).build()
            chain.proceed(updatedRequest)
        }
        .build()

    fun getAvailableApps(): List<TrueNasAppDetailsResponse> =
        executeGetRequest("app") {
            it.body?.string()?.let(jsonParser::decodeFromString) ?: emptyList()
        }

    private fun <T> executeGetRequest(endpoint: String, handler: (Response) -> T): T {
        val request = Request.Builder().url("$baseUrl/api/v2.0/$endpoint").get().build()
        return client.newCall(request).execute().use(handler)
    }
}
