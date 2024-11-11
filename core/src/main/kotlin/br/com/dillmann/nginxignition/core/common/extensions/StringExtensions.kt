package br.com.dillmann.nginxignition.core.common.extensions

import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

@OptIn(ExperimentalEncodingApi::class)
internal fun String.decodeBase64() = Base64.decode(this)
