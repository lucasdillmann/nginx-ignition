val koinVersion: String by project
val mapStructVersion: String by project

plugins {
    kotlin("kapt")
    id("io.ktor.plugin") version "3.0.1"
    id("org.jetbrains.kotlin.plugin.serialization") version "2.0.21"
}

application {
    mainClass.set("io.ktor.server.netty.EngineMain")

    val isDevelopment: Boolean = project.ext.has("development")
    applicationDefaultJvmArgs = listOf("-Dio.ktor.development=$isDevelopment")
}

dependencies {
    implementation(project(":core"))
    implementation(project(":database"))
    implementation(project(":custom-certificate"))
    implementation(project(":self-signed-certificate"))
    implementation(project(":acme-certificate"))
    implementation(project(":frontend"))
    implementation("io.ktor:ktor-server-core-jvm")
    implementation("io.ktor:ktor-server-auth-jvm")
    implementation("io.ktor:ktor-server-auth-jwt-jvm")
    implementation("io.ktor:ktor-server-host-common-jvm")
    implementation("io.ktor:ktor-server-caching-headers-jvm")
    implementation("io.ktor:ktor-server-compression-jvm")
    implementation("io.ktor:ktor-server-content-negotiation-jvm")
    implementation("io.ktor:ktor-serialization-kotlinx-json-jvm")
    implementation("io.insert-koin:koin-ktor:$koinVersion")
    implementation("io.insert-koin:koin-logger-slf4j:$koinVersion")
    implementation("io.ktor:ktor-server-netty-jvm")
    implementation("io.ktor:ktor-server-config-yaml")
    implementation("ch.qos.logback:logback-classic:1.5.12")
    implementation("org.mapstruct:mapstruct:$mapStructVersion")
    kapt("org.mapstruct:mapstruct-processor:$mapStructVersion")
}

ktor {
    docker {
        customBaseImage = "eclipse-temurin:21-jre-alpine"
        jreVersion = JavaVersion.VERSION_21
        localImageName = "dillmann/nginx-ignition"
        imageTag = rootProject.version.toString()
    }
}
