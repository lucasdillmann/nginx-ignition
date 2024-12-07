package br.com.dillmann.nginxignition.certificate.acme.dns.cloudflare

import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonNamingStrategy
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody

internal class CloudflareApiClient(private val apiToken: String) {
    @OptIn(ExperimentalSerializationApi::class)
    private val serializer = Json {
        ignoreUnknownKeys = true
        namingStrategy = JsonNamingStrategy.SnakeCase
        explicitNulls = false
    }

    private val client = OkHttpClient
        .Builder()
        .addInterceptor { chain ->
            val updatedRequest =
                chain.request().newBuilder().addHeader("Authorization", "Bearer $apiToken").build()
            chain.proceed(updatedRequest)
        }
        .build()

    fun listZones(): List<CloudflareApi.Zone> =
        fetchAll(atPath("zones").get().build())

    fun listRecords(zoneId: String): List<CloudflareApi.DnsRecord> =
        fetchAll(atPath("zones/$zoneId/dns_records").get().build())

    fun batchRecordChange(zoneId: String, batch: CloudflareApi.BatchDnsRecordChange) {
        val payload = serializer.encodeToString(batch).toRequestBody("application/json".toMediaType())
        val request = atPath("zones/$zoneId/dns_records/batch").post(payload).build()
        client.newCall(request).execute().use { response ->
            require(response.isSuccessful)
        }
    }

    private fun atPath(path: String) =
        Request.Builder().url("https://api.cloudflare.com/client/v4/$path")

    private inline fun <reified T> fetchAll(template: Request): List<T> {
        val output = mutableListOf<T>()
        var pageNumber = 1

        do {
            val url = template.url.newBuilder().addQueryParameter("page", pageNumber++.toString()).build()
            val request = template.newBuilder().url(url).build()
            val page = client.newCall(request).execute().use { response ->
                serializer.decodeFromString<CloudflareApi.Page<T>>(response.body!!.string())
            }

            output += page.result
        } while (page.success && page.resultInfo.hasMore)

        return output
    }
}
