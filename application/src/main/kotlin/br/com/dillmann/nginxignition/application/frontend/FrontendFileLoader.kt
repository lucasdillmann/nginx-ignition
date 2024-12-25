package br.com.dillmann.nginxignition.application.frontend

import java.net.URL

internal object FrontendFileLoader {
    private const val PATH_PREFIX = "nginx-ignition/frontend"

    fun load(path: String): URL? =
        @Suppress("StringLiteralDuplication")
        with(normalizePath("$PATH_PREFIX$path")) {
            loadUsingClassLoader(this)
                ?: loadUsingParentClassLoader(this)
                ?: loadUsingThreadLoader(this)
                ?: loadUsingClassLoader("/$this")
                ?: loadUsingParentClassLoader("/$this")
                ?: loadUsingThreadLoader("/$this")
        }

    private fun loadUsingClassLoader(path: String): URL? =
        FrontendFileLoader::class.java.classLoader?.getResource(path)

    private fun loadUsingParentClassLoader(path: String): URL? =
        FrontendFileLoader::class.java.classLoader?.parent?.getResource(path)

    private fun loadUsingThreadLoader(path: String): URL? =
        Thread.currentThread().contextClassLoader?.getResource(path)

    private tailrec fun normalizePath(path: String): String =
        if (path.contains("//")) normalizePath(path.replace("//", "/"))
        else path
}
