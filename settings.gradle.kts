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

project(":custom-certificate").projectDir = file("certificate/custom")
project(":lets-encrypt-certificate").projectDir = file("certificate/lets-encrypt")
project(":self-signed-certificate").projectDir = file("certificate/self-signed")
