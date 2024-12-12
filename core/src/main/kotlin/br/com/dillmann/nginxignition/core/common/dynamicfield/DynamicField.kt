package br.com.dillmann.nginxignition.core.common.dynamicfield

data class DynamicField(
    val id: String,
    val priority: Int,
    val description: String,
    val required: Boolean,
    val type: Type,
    val enumOptions: List<EnumOption> = emptyList(),
    val sensitive: Boolean = false,
    val condition: Condition? = null,
    val helpText: String? = null,
    val defaultValue: String? = null,
) {
    data class EnumOption(
        val id: String,
        val description: String,
    )

    data class Condition(
        val parentField: String,
        val value: String,
    )

    enum class Type {
        SINGLE_LINE_TEXT,
        MULTI_LINE_TEXT,
        URL,
        EMAIL,
        BOOLEAN,
        ENUM,
        FILE,
    }
}
