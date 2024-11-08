package br.com.dillmann.nginxignition.application.controller.certificate.model

import kotlinx.serialization.Serializable

@Serializable
data class AvailableProviderResponse(
    val uniqueId: String,
    val name: String,
    val dynamicFields: List<DynamicField>,
) {
    @Serializable
    data class DynamicField(
        val id: String,
        val description: String,
        val required: Boolean,
        val type: Type,
        val enumOptions: List<EnumOption> = emptyList(),
        val condition: Condition? = null,
    ) {
        @Serializable
        data class EnumOption(
            val id: String,
            val description: String,
        )

        @Serializable
        data class Condition(
            val parentField: String,
            val value: String,
        )

        enum class Type {
            SINGLE_LINE_TEXT,
            MULTI_LINE_TEXT,
            EMAIL,
            BOOLEAN,
            ENUM,
            FILE,
        }
    }
}
