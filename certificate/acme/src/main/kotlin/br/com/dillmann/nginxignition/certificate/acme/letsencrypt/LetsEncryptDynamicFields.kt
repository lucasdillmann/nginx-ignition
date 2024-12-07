package br.com.dillmann.nginxignition.certificate.acme.letsencrypt

import br.com.dillmann.nginxignition.certificate.acme.dns.aws.Route53DnsProvider
import br.com.dillmann.nginxignition.certificate.acme.dns.cloudflare.CloudflareDnsProvider
import br.com.dillmann.nginxignition.certificate.acme.dns.google.GoogleCloudDnsProvider
import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object LetsEncryptDynamicFields {
    val TERMS_OF_SERVICE = DynamicField(
        id = "acceptTheTermsOfService",
        priority = 99,
        description = "Terms of service",
        helpText = "I agree to the Let's Encrypt terms of service available at theirs website",
        required = true,
        type = DynamicField.Type.BOOLEAN,
    )

    val EMAIL_ADDRESS = DynamicField(
        id = "emailAddress",
        priority = 0,
        description = "E-mail address",
        required = true,
        type = DynamicField.Type.EMAIL,
    )

    val DNS_PROVIDER = DynamicField(
        id = "challengeDnsProvider",
        priority = 1,
        description = "DNS provider (for the DNS challenge)",
        required = true,
        type = DynamicField.Type.ENUM,
        enumOptions = listOf(
            DynamicField.EnumOption(Route53DnsProvider.ID, "AWS Route53"),
            DynamicField.EnumOption(CloudflareDnsProvider.ID, "Cloudflare DNS"),
            DynamicField.EnumOption(GoogleCloudDnsProvider.ID, "Google Cloud DNS"),
        ),
    )

    val AWS_ACCESS_KEY = DynamicField(
        id = "awsAccessKey",
        priority = 2,
        description = "AWS access key (for the DNS challenge)",
        required = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
        condition = DynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = Route53DnsProvider.ID,
        ),
    )

    val AWS_SECRET_KEY = DynamicField(
        id = "awsSecretKey",
        priority = 3,
        description = "AWS secret key",
        required = true,
        sensitive = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
        condition = DynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = Route53DnsProvider.ID,
        ),
    )

    val CLOUDFLARE_API_TOKEN = DynamicField(
        id = "cloudflareApiToken",
        priority = 2,
        description = "Cloudflare API token (for the DNS challenge)",
        required = true,
        sensitive = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
        condition = DynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = CloudflareDnsProvider.ID,
        ),
    )

    val GOOGLE_CLOUD_PRIVATE_KEY = DynamicField(
        id = "googleCloudPrivateKey",
        priority = 2,
        description = "Service account private key JSON",
        required = true,
        sensitive = true,
        type = DynamicField.Type.MULTI_LINE_TEXT,
        condition = DynamicField.Condition(
            parentField = DNS_PROVIDER.id,
            value = GoogleCloudDnsProvider.ID,
        ),
    )
}
