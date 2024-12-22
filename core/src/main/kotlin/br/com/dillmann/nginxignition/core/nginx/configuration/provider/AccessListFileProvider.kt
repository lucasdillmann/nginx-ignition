package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.accesslist.AccessList
import br.com.dillmann.nginxignition.core.accesslist.AccessListRepository
import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider
import org.apache.commons.codec.digest.Md5Crypt

internal class AccessListFileProvider(
    private val accessListRepository: AccessListRepository,
): NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        accessListRepository.findAll().flatMap { build(it, basePath) }

    private fun build(accessList: AccessList, basePath: String): List<NginxConfigurationFileProvider.Output> =
        listOfNotNull(
            buildConfFile(accessList, basePath),
            buildHtpasswdFile(accessList),
        )

    private fun buildConfFile(accessList: AccessList, basePath: String): NginxConfigurationFileProvider.Output {
        val entriesContents = accessList
            .entries
            .joinToString(separator = "\n") { entry ->
                entry
                    .sourceAddresses
                    .joinToString(separator = "\n") {
                        "${entry.outcome.toNginxOperation()} $it;"
                    }
            }

        val usernamePasswordContents =
            if (accessList.credentials.isEmpty()) ""
            else """
                auth_basic "${accessList.realm ?: "Login required"}";
                auth_basic_user_file $basePath/config/access-list-${accessList.id}.htpasswd;
            """.trimIndent()

        val satisfyContents =
            if (accessList.credentials.isEmpty() || accessList.entries.isEmpty()) "satisfy any;"
            else "satisfy ${if(accessList.satisfyAll) "all" else "any"};"

        val forwardHeadersContents =
            if (accessList.forwardAuthenticationHeader) ""
            else """
                proxy_set_header Authorization "";
            """.trimIndent()

        val contents = """
            $satisfyContents
            $entriesContents
            ${accessList.defaultOutcome.toNginxOperation()} all;
            $usernamePasswordContents
            $forwardHeadersContents
        """.trimIndent()

        return NginxConfigurationFileProvider.Output(
            name = "access-list-${accessList.id}.conf",
            contents = contents,
        )
    }

    private fun buildHtpasswdFile(accessList: AccessList): NginxConfigurationFileProvider.Output? {
        if (accessList.credentials.isEmpty())
            return null

        val contents = accessList
            .credentials
            .joinToString(separator = "\n") { (username, password) ->
                val hash = Md5Crypt.apr1Crypt(password)
                "$username:$hash"
            }

        return NginxConfigurationFileProvider.Output(
            name = "access-list-${accessList.id}.htpasswd",
            contents = contents,
        )
    }

    private fun AccessList.Outcome.toNginxOperation() =
        when (this) {
            AccessList.Outcome.DENY -> "deny"
            AccessList.Outcome.ALLOW -> "allow"
        }
}
