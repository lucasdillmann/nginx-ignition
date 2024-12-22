package br.com.dillmann.nginxignition.core.accesslist

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyValidator
import br.com.dillmann.nginxignition.core.common.validation.ErrorCreator
import inet.ipaddr.IPAddressString

internal class AccessListValidator : ConsistencyValidator() {
    private companion object {
        private const val VALUE_MISSING_MESSAGE = "A value is required"
    }

    suspend fun validate(accessList: AccessList) {
        withValidationScope { addError ->
            if (accessList.name.isBlank())
                addError("name", VALUE_MISSING_MESSAGE)

            val previousUsernames = mutableListOf<String>()
            accessList.credentials.forEachIndexed { index, credentials ->
                validateCredentials(index, credentials, previousUsernames, addError)
            }

            val previousPriorities = mutableSetOf<Int>()
            accessList.entries.forEachIndexed { index, entrySet ->
                validateEntrySet(index, entrySet, previousPriorities, addError)
            }
        }
    }

    private fun validateEntrySet(
        index: Int,
        entrySet: AccessList.EntrySet,
        knownPriorities: MutableSet<Int>,
        addError: ErrorCreator,
    ) {
        if (!knownPriorities.add(entrySet.priority))
            addError("entries[$index].priority", "Value is duplicated")

        if (entrySet.priority < 0)
            addError("entries[$index].priority", "Value must be 0 or greater")

        if (entrySet.sourceAddresses.isEmpty())
            addError("entries[$index].sourceAddresses", VALUE_MISSING_MESSAGE)

        entrySet.sourceAddresses.forEach { sourceAddress ->
            if (!isAValidIpAddressOrRange(sourceAddress))
                addError(
                    "entries[$index].sourceAddresses",
                    "Address \"$sourceAddress\" is not a valid IPv4 or IPv6 address or range",
                )
        }
    }

    private fun isAValidIpAddressOrRange(input: String): Boolean =
        try {
            IPAddressString(input).isValid
        } catch (_: Exception) {
            false
        }

    private fun validateCredentials(
        index: Int,
        credentials: AccessList.Credentials,
        previousUsernames: MutableList<String>,
        addError: ErrorCreator,
    ) {
        if (credentials.username.isBlank())
            addError("credentials[$index].username", VALUE_MISSING_MESSAGE)
        else if (!previousUsernames.add(credentials.username))
            addError("credentials[$index].username", "Value is duplicated")
    }
}
