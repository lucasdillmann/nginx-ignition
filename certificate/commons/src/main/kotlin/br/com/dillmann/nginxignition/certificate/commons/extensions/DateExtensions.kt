package br.com.dillmann.nginxignition.certificate.commons.extensions

import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.*

fun Date.toOffsetDateTime(): OffsetDateTime = toInstant().atOffset(ZoneOffset.UTC)
