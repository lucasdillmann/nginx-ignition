plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "0.8.0"
}

rootProject.name = "nginx-side-wheel"
include("core", "database", "application", "third-party", "frontend")

