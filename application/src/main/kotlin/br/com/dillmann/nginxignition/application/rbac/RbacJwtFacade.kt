package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import br.com.dillmann.nginxignition.core.user.command.GetUserStatusCommand
import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import com.auth0.jwt.interfaces.DecodedJWT
import java.security.SecureRandom
import java.time.Instant
import java.util.UUID

class RbacJwtFacade(
    configurationProvider: ConfigurationProvider,
    private val getUserStatusCommand: GetUserStatusCommand,
    private val getUserCommand: GetUserCommand,
) {
    companion object {
        const val UNIQUE_IDENTIFIER = "nginx-ignition"
        private val REVOKED_IDS = mutableSetOf<String>()
        private const val EXPECTED_JWT_SECRET_SIZE_CHARS = 64
        private const val EXPECTED_JWT_SECRET_SIZE_BYTES = 512
    }

    private val logger = logger<RbacJwtFacade>()
    private val configurationProvider = configurationProvider.withPrefix("nginx-ignition.security.jwt")
    private val jwtSecret = initializeSecret()
    private val verifier = JWT
        .require(Algorithm.HMAC512(jwtSecret))
        .withAudience(UNIQUE_IDENTIFIER)
        .withIssuer(UNIQUE_IDENTIFIER)
        .build()

    suspend fun checkCredentials(call: ApiCall): Subject? {
        val token = call.jwtToken() ?: return null
        val jwt = checkCredentials(token) ?: return null
        return Subject(
            tokenId = jwt.id,
            userId = UUID.fromString(jwt.subject),
        )
    }

    suspend fun checkCredentials(token: String): DecodedJWT? {
        val jwt = runCatching { verifier.verify(token) }.getOrNull() ?: return null
        val userId = runCatching { jwt.subject.let(UUID::fromString) }.getOrNull()

        return when {
            UNIQUE_IDENTIFIER !in jwt.audience -> null
            UNIQUE_IDENTIFIER != jwt.issuer -> null
            jwt.id in REVOKED_IDS -> null
            userId == null -> null
            !getUserStatusCommand.isEnabled(userId) -> {
                revokeCredentials(jwt.id)
                null
            }
            else -> jwt
        }
    }

    fun revokeCredentials(tokenId: String) {
        REVOKED_IDS += tokenId
    }

    fun buildToken(user: User): String {
        val ttlSeconds = configurationProvider.get("ttl-seconds").toLong()
        val clockSkewSeconds = configurationProvider.get("clock-skew-seconds").toLong()

        return JWT
            .create()
            .withAudience(UNIQUE_IDENTIFIER)
            .withIssuer(UNIQUE_IDENTIFIER)
            .withJWTId(UUID.randomUUID().toString())
            .withNotBefore(Instant.now().minusSeconds(clockSkewSeconds))
            .withExpiresAt(Instant.now().plusSeconds(ttlSeconds).plusSeconds(clockSkewSeconds))
            .withIssuedAt(Instant.now())
            .withSubject(user.id.toString())
            .withClaim("username", user.username)
            .withClaim("role", user.role.name)
            .sign(Algorithm.HMAC512(jwtSecret))
    }

    suspend fun refreshToken(credentials: DecodedJWT): String? {
        val windowSize = configurationProvider.get("renew-window-seconds").toLong()
        val expiration = credentials.expiresAt?.toInstant() ?: return null
        val renewWindow = expiration.minusSeconds(windowSize)..expiration
        if (Instant.now() !in renewWindow) return null

        val user = getUserCommand.getById(credentials.subject.let(UUID::fromString)) ?: return null
        return buildToken(user)
    }

    private fun initializeSecret(): ByteArray {
        val secret = configurationProvider.get("secret")
        if (secret.isNotBlank()) {
            if(secret.length != EXPECTED_JWT_SECRET_SIZE_CHARS) {
                val message =
                    "JWT secret should be 64 characters long (512 bytes) but is ${secret.length} characters long"
                logger.error(message)
                throw IllegalArgumentException(message)
            }

            return secret.toByteArray(charset = Charsets.UTF_8)
        }

        logger.warn(
            "Application was initialized without a JWT secret and a random one will be generated. This will lead " +
                "to users being logged-out every time the app restarts or they hit a different instance. Please " +
                "refer to the documentation in order to provide a custom secret.",
        )

        return ByteArray(EXPECTED_JWT_SECRET_SIZE_BYTES).also(SecureRandom()::nextBytes)
    }
}
