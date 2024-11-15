package br.com.dillmann.nginxignition.core.common

object GlobalConstants {
    val EMAIL_PATTERN = "^[A-Za-z0-9+_.-]+@[A-Za-z0-9.-]+\$".toPattern()
    val TLD_PATTERN = "(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]".toPattern()
}
