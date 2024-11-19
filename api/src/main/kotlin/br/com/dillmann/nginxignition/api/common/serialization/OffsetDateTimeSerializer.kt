package br.com.dillmann.nginxignition.api.common.serialization

import kotlinx.serialization.KSerializer
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import java.time.OffsetDateTime

typealias OffsetDateTimeString = @Serializable(with = OffsetDateTimeSerializer::class) OffsetDateTime

class OffsetDateTimeSerializer: KSerializer<OffsetDateTime> {
    override val descriptor = PrimitiveSerialDescriptor("offsetDateTime", PrimitiveKind.STRING)

    override fun deserialize(decoder: Decoder): OffsetDateTime =
        OffsetDateTime.parse(decoder.decodeString())

    override fun serialize(encoder: Encoder, value: OffsetDateTime) {
        encoder.encodeString(value.toString())
    }
}
