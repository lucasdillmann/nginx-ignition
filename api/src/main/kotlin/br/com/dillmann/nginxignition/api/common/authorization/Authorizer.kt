package br.com.dillmann.nginxignition.api.common.authorization

interface Authorizer {
    suspend fun revoke(subject: Subject)
    suspend fun buildToken(subject: Subject): String
}
