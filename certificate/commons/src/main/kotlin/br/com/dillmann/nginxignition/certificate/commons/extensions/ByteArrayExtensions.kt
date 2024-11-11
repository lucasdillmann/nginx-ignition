package br.com.dillmann.nginxignition.certificate.commons.extensions

import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

@OptIn(ExperimentalEncodingApi::class)
fun ByteArray.encodeBase64() = Base64.encode(this)
