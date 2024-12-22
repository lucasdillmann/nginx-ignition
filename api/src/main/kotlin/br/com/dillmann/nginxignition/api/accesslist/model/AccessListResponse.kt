package br.com.dillmann.nginxignition.api.accesslist.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import br.com.dillmann.nginxignition.core.accesslist.AccessList
import kotlinx.serialization.Serializable

@Serializable
data class AccessListResponse(
    val id: UuidString,
    val name: String,
    val realm: String?,
    val satisfyAll: Boolean,
    val defaultOutcome: AccessList.Outcome,
    val entries: List<EntrySet>,
    val forwardAuthenticationHeader: Boolean,
    val credentials: List<Credentials>,
) {
    @Serializable
    data class EntrySet(
        val priority: Int,
        val outcome: AccessList.Outcome,
        val sourceAddresses: List<String>,
    )

    @Serializable
    data class Credentials(
        val username: String,
        val password: String,
    )
}
