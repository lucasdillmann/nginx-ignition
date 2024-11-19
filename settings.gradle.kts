plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "0.8.0"
}

rootProject.name = "nginx-ignition"
include(
    "core",
    "database",
    "application",
    "frontend",
    "api",
    "certificate-commons",
    "acme-certificate",
    "custom-certificate",
    "self-signed-certificate",
)

project(":custom-certificate").projectDir = file("certificate/custom")
project(":certificate-commons").projectDir = file("certificate/commons")
project(":acme-certificate").projectDir = file("certificate/acme")
project(":self-signed-certificate").projectDir = file("certificate/self-signed")
