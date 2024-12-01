package br.com.dillmann.nginxignition.integration.truenas

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

internal object TrueNasDynamicFields {
    val URL = DynamicField(
        id = "url",
        description = "URL",
        priority = 1,
        required = true,
        helpText = "The URL where your NAS is accessible, like http://192.168.0.2 or https://nas.yourdomain.com",
        type = DynamicField.Type.SINGLE_LINE_TEXT,
    )

    val USERNAME = DynamicField(
        id = "username",
        description = "Username",
        priority = 2,
        required = true,
        sensitive = false,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
    )

    val PASSWORD = DynamicField(
        id = "password",
        description = "Password",
        priority = 3,
        required = true,
        sensitive = true,
        type = DynamicField.Type.SINGLE_LINE_TEXT,
    )
}
