package br.com.dillmann.nginxignition.core.common.log

import org.slf4j.Logger
import org.slf4j.LoggerFactory

inline fun <reified T : Any> logger(): Logger =
    LoggerFactory.getLogger(T::class.java)

fun logger(name: String): Logger =
    LoggerFactory.getLogger(name)
