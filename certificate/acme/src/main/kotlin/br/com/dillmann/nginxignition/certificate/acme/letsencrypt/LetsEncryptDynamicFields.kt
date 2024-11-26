package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

import br.com.dillmann.nginxignition.core.certificate.provider.CertificateProviderDynamicField

internal object LetsEncryptDynamicFields {
    val TERMS_OF_SERVICE = CertificateProviderDynamicField(
        id = "acceptTheTermsOfService",
        priority = 99,
        description = "Terms of service",
        helpText = "I agree to the Let's Encrypt terms of service available at theirs site",
        required = true,
        type = CertificateProviderDynamicField.Type.BOOLEAN,
    )

    val EMAIL_ADDRESS = CertificateProviderDynamicField(
        id = "emailAddress",
        priority = 0,
        description = "E-mail address",
        required = true,
        type = CertificateProviderDynamicField.Type.EMAIL,
    )

    val DNS_PROVIDER = CertificateProviderDynamicField(
        id = "challengeDnsProvider",
        priority = 1,
        description = "DNS provider (for the DNS challenge)",
        required = true,
        type = CertificateProviderDynamicField.Type.ENUM,
        enumOptions = listOf(
            CertificateProviderDynamicField.EnumOption("AWS_ROUTE53", "AWS Route53"),
        ),
    )

    val AWS_ACCESS_KEY = CertificateProviderDynamicField(
        id = "awsAccessKey",
        priority = 2,
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
        priority = 3,
        description = "AWS secret key",
        required = true,
        sensitive = true,
        type = CertificateProviderDynamicField.Type.SINGLE_LINE_TEXT,
        condition = CertificateProviderDynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = "AWS_ROUTE53",
        ),
    )
}
