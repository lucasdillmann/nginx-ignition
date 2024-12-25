val koinVersion: String by project
val mapStructVersion: String by project
val logbackVersion: String by project
val kotlinSerializationVersion: String by project
val snakeYamlVersion: String by project
val jwtVersion: String by project
val apacheTikaVersion: String by project

plugins {
    application
    kotlin("kapt")
    id("org.jetbrains.kotlin.plugin.serialization") version "2.1.0"
    id("com.github.johnrengelman.shadow") version "8.1.1"
}

application {
    mainClass = "br.com.dillmann.nginxignition.application.MainKt"
}

dependencies {
    implementation(project(":core"))
    implementation(project(":database"))
    implementation(project(":custom-certificate"))
    implementation(project(":self-signed-certificate"))
    implementation(project(":acme-certificate"))
    implementation(project(":truenas-integration"))
    implementation(project(":docker-integration"))
    implementation(project(":frontend"))
    implementation(project(":api"))
    implementation("io.insert-koin:koin-core-jvm:$koinVersion")
    implementation("io.insert-koin:koin-logger-slf4j:$koinVersion")
    implementation("ch.qos.logback:logback-classic:$logbackVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json-jvm:$kotlinSerializationVersion")
    implementation("org.yaml:snakeyaml:$snakeYamlVersion")
    implementation("com.auth0:java-jwt:$jwtVersion")
    implementation("org.apache.tika:tika-core:$apacheTikaVersion")
}

tasks {
    shadowJar {
        mergeServiceFiles()
    }
}
