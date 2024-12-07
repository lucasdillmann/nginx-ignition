package br.com.dillmann.nginxignition.certificate.acme.dns.cloudflare

import kotlinx.serialization.Serializable

internal class CloudflareApi {
    @Serializable
    data class Page<T>(
        val success: Boolean,
        val resultInfo: PageInfo,
        val result: List<T>,
    )

    @Serializable
    data class PageInfo(
        val totalCount: Long,
        val page: Long,
        val count: Long,
        val perPage: Long,
    ) {
        val hasMore = page * perPage < totalCount
    }

    @Serializable
    data class Zone (
        val id: String,
        val name: String,
    )

    @Serializable
    data class DnsRecord(
        val id: String? = null,
        val name: String,
        val content: String,
        val type: String,
    )

    @Serializable
    data class BatchDnsRecordChange(
        val posts: List<DnsRecord>,
        val puts: List<DnsRecord>,
    )
}
