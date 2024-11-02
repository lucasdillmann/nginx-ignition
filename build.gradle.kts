val koinVersion: String by project

plugins {
    kotlin("jvm") version "2.0.21"
}

group = "br.com.dillmann.nginxsidewheel"
version = "1.0.0"

allprojects {
    repositories {
        mavenCentral()
    }
}

subprojects {
    apply(plugin = "kotlin")

    dependencies {
        implementation("io.insert-koin:koin-core-jvm:$koinVersion")
    }
}
