package br.com.dillmann.nginxignition.integration.truenas

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object TrueNasDynamicFields {
    val URL = DynamicField(
        id = "url",
        description = "URL",
        priority = 1,
        required = true,
        helpText = "The URL/IP where your NAS is accessible, like http://192.168.0.2 or https://nas.yourdomain.com",
        type = DynamicField.Type.URL,
    )

    val PROXY_URL = DynamicField(
        id = "proxyUrl",
        description = "Apps URL",
        priority = 2,
        required = false,
        helpText = "The URL/IP to be used when proxying a request to a TrueNAS app. Use this if the apps are " +
            "exposed in another address that isn't the same as the main URL above.",
        type = DynamicField.Type.URL,
    )

    val USERNAME = DynamicField(
        id = "username",
        description = "Username",
        priority = 3,
        required = true,
        sensitive = false,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
    )

    val PASSWORD = DynamicField(
        id = "password",
        description = "Password",
        priority = 4,
        required = true,
        sensitive = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
    )
}
