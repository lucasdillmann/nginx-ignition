package br.com.dillmann.nginxignition.api.common.jsonobject

import kotlinx.serialization.json.*

fun Map<String, Any?>.toJsonObject(): JsonObject {
    val contents = map { (key, value) -> key to value.wrap() }.toMap()
    return JsonObject(contents)
}

fun JsonObject.toUnwrappedMap(): Map<String, Any?> =
    entries.associate { (key, rawValue) ->
        val value =
            if (rawValue is JsonObject) rawValue.toUnwrappedMap()
            else rawValue.jsonPrimitive.unwrapValue()

        key to value
    }

private fun JsonPrimitive.unwrapValue(): Any? =
    booleanOrNull
        ?: longOrNull?.toBigInteger()
        ?: doubleOrNull?.toBigDecimal()
        ?: contentOrNull

private fun Any?.wrap(): JsonPrimitive =
    when (this) {
        null -> JsonNull
        is Number -> JsonPrimitive(this)
        is Boolean -> JsonPrimitive(this)
        is String -> JsonPrimitive(this)
        else -> JsonPrimitive(toString())
    }
