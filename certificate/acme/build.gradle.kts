val acme4jVersion: String by project
val kotlinSerializationVersion: String by project
val awsSdkVersion: String by project
val azureSdkDnsVersion: String by project
val azureSdkIdentityVersion: String by project
val azureSdkVaultVersion: String by project
val googleCloudSdkVersion: String by project
val okHttpVersion: String by project

plugins {
    id("org.jetbrains.kotlin.plugin.serialization") version "2.1.0"
}

dependencies {
    implementation(project(":core"))
    implementation(project(":certificate-commons"))
    implementation("org.shredzone.acme4j:acme4j-client:$acme4jVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$kotlinSerializationVersion")
    implementation("software.amazon.awssdk:route53:$awsSdkVersion")
    implementation("com.squareup.okhttp3:okhttp:$okHttpVersion")
    implementation("com.google.cloud:google-cloud-dns:$googleCloudSdkVersion")
    implementation("com.azure.resourcemanager:azure-resourcemanager-dns:$azureSdkDnsVersion")
    implementation("com.azure:azure-identity:$azureSdkIdentityVersion")
    implementation("com.azure:azure-security-keyvault-secrets:$azureSdkVaultVersion")
}
