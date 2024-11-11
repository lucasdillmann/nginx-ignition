package br.com.dillmann.nginxignition.certificate.custom.extensions

import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.*

internal fun Date.toOffsetDateTime(): OffsetDateTime = toInstant().atOffset(ZoneOffset.UTC)
