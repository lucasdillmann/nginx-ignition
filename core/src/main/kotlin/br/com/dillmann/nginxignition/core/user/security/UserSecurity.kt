package br.com.dillmann.nginxignition.core.user.security

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import java.security.MessageDigest
import java.security.SecureRandom
import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

internal class UserSecurity(configurationProvider: ConfigurationProvider) {
    data class Output(
        val hash: String,
        val salt: String,
    )

    private val configurationProvider =
        configurationProvider.withPrefix("nginx-ignition.security.user-password-hashing")
    private val randomGenerator =
        SecureRandom()

    @OptIn(ExperimentalEncodingApi::class)
    fun hash(password: String): Output {
        val saltSize = configurationProvider.get("salt-size").toInt()
        val passwordBytes = password.encodeToByteArray()
        val saltBytes = ByteArray(saltSize)
        randomGenerator.nextBytes(saltBytes)

        val hashBytes = hash(passwordBytes, saltBytes)
        return Output(
            hash = Base64.encode(hashBytes),
            salt = Base64.encode(saltBytes),
        )
    }

    @OptIn(ExperimentalEncodingApi::class)
    fun check(
        providedPassword: String,
        storedHash: String,
        storedSalt: String,
    ): Boolean {
        val passwordBytes = providedPassword.encodeToByteArray()
        val saltBytes = Base64.decode(storedSalt)
        val hashBytes = hash(passwordBytes, saltBytes)
        return Base64.encode(hashBytes) == storedHash
    }

    private fun hash(password: ByteArray, salt: ByteArray): ByteArray {
        val hashAlgorithm = configurationProvider.get("algorithm")
        val hashIterations = configurationProvider.get("iterations").toInt()
        val hashGenerator = MessageDigest.getInstance(hashAlgorithm)
        var hashOutput = password + salt
        repeat(hashIterations) {
            hashOutput = hashGenerator.digest(hashOutput)
        }

        return hashOutput
    }
}
