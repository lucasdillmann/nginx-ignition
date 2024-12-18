package br.com.dillmann.nginxignition.core.accesslist

import java.util.UUID

data class AccessList(
    val id: UUID,
    val name: String,
    val realm: String?,
    val defaultOutcome: Outcome,
    val entries: List<EntrySet>,
    val forwardAuthenticationHeader: Boolean,
    val credentials: List<Credentials>,
) {
    enum class Outcome {
        DENY,
        ALLOW,
    }

    data class EntrySet(
        val priority: Int,
        val outcome: Outcome,
        val sourceAddresses: List<String>,
    )

    data class Credentials(
        val username: String,
        val password: String,
    )
}
