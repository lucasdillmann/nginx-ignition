package br.com.dillmann.nginxignition.certificate.letsencrypt

import kotlinx.serialization.Serializable

@Serializable
internal data class LetsEncryptMetadata(
    val userMail: String,
    val userPrivateKey: String,
    val userPublicKey: String,
)
