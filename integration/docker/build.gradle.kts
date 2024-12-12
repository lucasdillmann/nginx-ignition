val dockerJavaVersion: String by project

plugins {
    id("org.jetbrains.kotlin.plugin.serialization") version "2.1.0"
}

dependencies {
    implementation(project(":core"))
    implementation("com.github.docker-java:docker-java:$dockerJavaVersion")
    implementation("com.github.docker-java:docker-java-transport-zerodep:$dockerJavaVersion")
}
