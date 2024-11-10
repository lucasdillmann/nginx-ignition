package br.com.dillmann.nginxignition.database.common.json

import kotlinx.serialization.json.*

internal fun Map<*, *>.toJsonObject(): JsonObject =
    entries
        .associate { (key, value) ->
            require(key is String) { "All keys should be a string" }

            val jsonElement: JsonElement =
                when (value) {
                    null -> JsonNull
                    is Map<*, *> -> value.toJsonObject()
                    is Number -> JsonPrimitive(value)
                    is Boolean -> JsonPrimitive(value)
                    is Enum<*> -> JsonPrimitive(value.name)
                    else -> JsonPrimitive(value.toString())
                }

            key to jsonElement
        }
        .toMap()
        .let(::JsonObject)

internal fun String.toJsonObject() =
    Json.parseToJsonElement(this).jsonObject

internal fun JsonObject.toPlainMap(): Map<String, Any?> =
    entries.associate { (key, rawValue) ->
        val value =
            if (rawValue is JsonObject) rawValue.toPlainMap()
            else rawValue.jsonPrimitive.unwrap()

        key to value
    }

private fun JsonPrimitive.unwrap(): Any? =
    booleanOrNull
        ?: longOrNull?.toBigInteger()
        ?: doubleOrNull?.toBigDecimal()
        ?: contentOrNull
