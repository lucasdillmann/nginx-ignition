package br.com.dillmann.nginxignition.integration.truenas.client

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import com.google.common.cache.CacheBuilder
import org.koin.mp.KoinPlatform.getKoin
import java.util.concurrent.TimeUnit

internal object TrueNasApiCache {
    private val configProvider =
        getKoin().get<ConfigurationProvider>()
    private val cacheTimeoutSeconds =
        configProvider.get("nginx-ignition.integration.truenas.api-cache-timeout-seconds").toLong()
    private val delegate = CacheBuilder
        .newBuilder()
        .expireAfterAccess(cacheTimeoutSeconds, TimeUnit.SECONDS)
        .build<String, Any>()

    @Suppress("UNCHECKED_CAST")
    fun <T> get(key: String, missProvider: () -> T): T =
        delegate[key, missProvider] as T
}
