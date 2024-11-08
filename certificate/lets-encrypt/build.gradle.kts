val acme4jVersion: String by project
val kotlinxSerializationJsonVersion: String by project
val awsSdkVersion: String by project

plugins {
    id("org.jetbrains.kotlin.plugin.serialization") version "2.0.21"
}

dependencies {
    implementation(project(":core"))
    implementation("org.shredzone.acme4j:acme4j-client:$acme4jVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$kotlinxSerializationJsonVersion")
    implementation("software.amazon.awssdk:route53:$awsSdkVersion")
}
