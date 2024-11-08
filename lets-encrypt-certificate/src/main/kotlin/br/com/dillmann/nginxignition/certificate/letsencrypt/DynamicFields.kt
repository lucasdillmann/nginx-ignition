package br.com.dillmann.nginxignition.certificate.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object DynamicFields {
    val TERMS_OF_SERVICE = CertificateProviderDynamicField(
        id = "TERMS_OF_SERVICE",
        description = "I agree to the Let's Encrypt terms of service available at theirs site",
        required = true,
        type = CertificateProviderDynamicField.Type.BOOLEAN,
    )

    val EMAIL_ADDRESS = CertificateProviderDynamicField(
        id = "EMAIL_ADDRESS",
        description = "E-mail address",
        required = true,
        type = CertificateProviderDynamicField.Type.EMAIL,
    )

    val DNS_PROVIDER = CertificateProviderDynamicField(
        id = "CHALLENGE_DNS_PROVIDER",
        description = "DNS provider (for the DNS challenge)",
        required = true,
        type = CertificateProviderDynamicField.Type.ENUM,
        enumOptions = listOf(
            CertificateProviderDynamicField.EnumOption("AWS_ROUTE53", "AWS Route53"),
        ),
    )

    val AWS_ACCESS_KEY = CertificateProviderDynamicField(
        id = "AWS_ACCESS_KEY",
        description = "AWS access key (for the Route 53 DNS challenge)",
        required = true,
        type = CertificateProviderDynamicField.Type.SINGLE_LINE_TEXT,
        condition = CertificateProviderDynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = "AWS_ROUTE53",
        ),
    )

    val AWS_SECRET_KEY = CertificateProviderDynamicField(
        id = "AWS_SECRET_KEY",
        description = "AWS secret key",
        required = true,
        type = CertificateProviderDynamicField.Type.SINGLE_LINE_TEXT,
        condition = CertificateProviderDynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = "AWS_ROUTE53",
        ),
    )
}
