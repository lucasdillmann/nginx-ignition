val kotlinSerializationVersion: String by project
val mapStructVersion: String by project

plugins {
    kotlin("kapt")
    id("org.jetbrains.kotlin.plugin.serialization") version "2.1.0"
}

dependencies {
    implementation(project(":core"))
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json-jvm:$kotlinSerializationVersion")
    implementation("org.mapstruct:mapstruct:$mapStructVersion")
    kapt("org.mapstruct:mapstruct-processor:$mapStructVersion")
}
