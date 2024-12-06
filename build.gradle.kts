val koinVersion: String by project
val slf4jVersion: String by project
val coroutinesVersion: String by project

plugins {
    kotlin("jvm") version "2.1.0"
    id("org.sonarqube") version "6.0.1.5171"
    id("io.gitlab.arturbosch.detekt") version "1.23.6"
}

allprojects {
    group = "br.com.dillmann.nginxignition"
    version = "1.1.0"

    repositories {
        mavenCentral()
    }

    sonarqube {
        properties {
            property("sonar.projectKey", "lucasdillmann_nginx-ignition")
            property("sonar.organization", "lucasdillmann")
        }
    }
}

subprojects {
    apply(plugin = "kotlin")
    apply(plugin = "org.sonarqube")
    apply(plugin = "io.gitlab.arturbosch.detekt")

    dependencies {
        implementation("io.insert-koin:koin-core-jvm:$koinVersion")
        implementation("org.slf4j:slf4j-api:$slf4jVersion")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core-jvm:$coroutinesVersion")
    }
}

kotlin {
    jvmToolchain(21)
}
