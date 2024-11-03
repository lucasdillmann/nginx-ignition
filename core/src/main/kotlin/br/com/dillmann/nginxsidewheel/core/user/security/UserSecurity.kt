package br.com.dillmann.nginxsidewheel.core.user.security

import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import java.security.MessageDigest
import java.security.SecureRandom
import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

internal class UserSecurity(configurationProvider: ConfigurationProvider) {
    data class Output(
        val hash: String,
        val salt: String,
    )

    private val hashAlgorithm =
        configurationProvider.get("nginx-sidewheel.security.user-password-hashing.algorithm")
    private val saltSize =
        configurationProvider.get("nginx-sidewheel.security.user-password-hashing.salt-size").toInt()
    private val hashIterations =
        configurationProvider.get("nginx-sidewheel.security.user-password-hashing.iterations").toInt()
    private val randomGenerator =
        SecureRandom()

    @OptIn(ExperimentalEncodingApi::class)
    fun hash(password: String): Output {
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
        val hashGenerator = MessageDigest.getInstance(hashAlgorithm)
        var hashOutput = password + salt
        repeat(hashIterations) {
            hashOutput = hashGenerator.digest(hashOutput)
        }

        return hashOutput
    }
}
