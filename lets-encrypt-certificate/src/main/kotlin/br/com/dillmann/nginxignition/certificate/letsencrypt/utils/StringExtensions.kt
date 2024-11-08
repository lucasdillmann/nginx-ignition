package br.com.dillmann.nginxignition.certificate.letsencrypt.utils

import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

@OptIn(ExperimentalEncodingApi::class)
fun String.decodeBase64() = Base64.decode(this)
