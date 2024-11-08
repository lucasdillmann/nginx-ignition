plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "0.8.0"
}

rootProject.name = "nginx-ignition"
include(
    "core",
    "database",
    "application",
    "frontend",
    "lets-encrypt-certificate",
    "custom-certificate",
    "self-signed-certificate",
)

