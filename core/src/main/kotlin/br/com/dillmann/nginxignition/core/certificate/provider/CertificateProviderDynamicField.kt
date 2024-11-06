package br.com.dillmann.nginxignition.core.certificate.provider

data class CertificateProviderDynamicField(
    val uniqueId: String,
    val description: String,
    val required: Boolean,
    val type: Type,
    val enumOptions: List<String> = emptyList(),
) {
    enum class Type {
        SINGLE_LINE_TEXT,
        MULTI_LINE_TEXT,
        EMAIL,
        BOOLEAN,
        ENUM,
    }
}
