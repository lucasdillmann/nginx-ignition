package br.com.dillmann.nginxignition.certificate.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object DynamicFields {
    val TERMS_OF_SERVICE = CertificateProviderDynamicField(
        id = "acceptTheTermsOfService",
        description = "I agree to the Let's Encrypt terms of service available at theirs site",
        required = true,
        type = CertificateProviderDynamicField.Type.BOOLEAN,
    )

    val EMAIL_ADDRESS = CertificateProviderDynamicField(
        id = "emailAddress",
        description = "E-mail address",
        required = true,
        type = CertificateProviderDynamicField.Type.EMAIL,
    )

    val DNS_PROVIDER = CertificateProviderDynamicField(
        id = "callengeDnsProvider",
        description = "DNS provider (for the DNS challenge)",
        required = true,
        type = CertificateProviderDynamicField.Type.ENUM,
        enumOptions = listOf(
            CertificateProviderDynamicField.EnumOption("AWS_ROUTE53", "AWS Route53"),
        ),
    )

    val AWS_ACCESS_KEY = CertificateProviderDynamicField(
        id = "awsAccessKey",
        description = "AWS access key (for the Route 53 DNS challenge)",
        required = true,
        type = CertificateProviderDynamicField.Type.SINGLE_LINE_TEXT,
        condition = CertificateProviderDynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = "AWS_ROUTE53",
        ),
    )

    val AWS_SECRET_KEY = CertificateProviderDynamicField(
        id = "awsSecretKey",
        description = "AWS secret key",
        required = true,
        type = CertificateProviderDynamicField.Type.SINGLE_LINE_TEXT,
        condition = CertificateProviderDynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = "AWS_ROUTE53",
        ),
    )
}
