val koinVersion: String by project

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
    }
}

kotlin {
    jvmToolchain(21)
}
