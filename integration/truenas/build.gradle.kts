val okHttpVersion: String by project
val kotlinSerializationVersion: String by project
val guavaVersion: String by project

plugins {
    id("org.jetbrains.kotlin.plugin.serialization") version "2.0.21"
}

dependencies {
    implementation(project(":core"))
    implementation("com.squareup.okhttp3:okhttp:$okHttpVersion")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json-jvm:$kotlinSerializationVersion")
    implementation("com.google.guava:guava:$guavaVersion")
}
