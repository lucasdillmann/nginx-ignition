package br.com.dillmann.nginxignition.database.common.database

enum class DatabaseType {
    POSTGRESQL,
    H2,
    ;

    companion object {
        fun fromJdbcUrl(url: String): DatabaseType =
            when {
                url.startsWith("jdbc:h2:") -> H2
                url.startsWith("jdbc:postgresql:") -> POSTGRESQL
                else -> error("Unsupported database type: $url")
            }
    }
}
