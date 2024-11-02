val koinVersion: String by project
val slf4jVersion: String by project
val coroutinesVersion: String by project

plugins {
    kotlin("jvm") version "2.0.21"
}

allprojects {
    group = "br.com.dillmann.nginxsidewheel"
    version = "1.0.0"

    repositories {
        mavenCentral()
    }
}

subprojects {
    apply(plugin = "kotlin")

    dependencies {
        implementation("io.insert-koin:koin-core-jvm:$koinVersion")
        implementation("org.slf4j:slf4j-api:$slf4jVersion")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core-jvm:$coroutinesVersion")
    }
}

kotlin {
    jvmToolchain(21)
}
